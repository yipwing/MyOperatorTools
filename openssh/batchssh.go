package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

type (
	ownHost struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
)

// TODO  next add json file read.

func createLogger() *log.Logger {
	file, err := os.OpenFile("./batch.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("failed to create log file.")
	}
	// defer file.Close()  // if close file handle. will not be write.
	return log.New(file, "", log.LstdFlags|log.Llongfile)
}

func readFile() []ownHost {
	file, err := os.Open("./hosts.json")
	if err != nil {
		log.Fatalln("failed to open json file")
	}
	defer file.Close()
	result := []ownHost{}
	buffer, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		log.Fatalln("cannot read hosts.json file or is not validate json format.")
	}
	return result
}

// func handler() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	// mission(ctx)
// 	test(ctx)
// }

// func getKeys() ssh.AuthMethod {
// 	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
// 		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
// 	}
// 	return nil
// }

func execute() error {
	logger := createLogger()
	remoteDir := flag.String("remoteDir", "/root", "remote directory")
	flag.Parse()
	if dErr := delete(); dErr != nil {
		return dErr
	}
	err := scopy(*remoteDir)
	if err != nil {
		return err
	}
	Hosts := readFile()
	// fmt.Printf("%d mission ready to go \n", len(Hosts))
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		ip := items.Host
		szPassword := items.Password
		config := &ssh.ClientConfig{
			User: items.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(szPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		client, err := ssh.Dial("tcp", ip, config)
		if err != nil {
			logger.Printf("%s %s Maybe is the password error or is the host not avaliable\n", ip, err.Error())
			continue
		}
		defer client.Close()
		sessionForUpdate, sUPErr := client.NewSession()
		if sUPErr != nil {
			logger.Printf("%s Failed to create session: %s\n" + ip + sUPErr.Error())
			continue
		}
		defer sessionForUpdate.Close()
		var outUpBuff, errUpBuff bytes.Buffer
		sessionForUpdate.Stdout = &outUpBuff
		sessionForUpdate.Stderr = &errUpBuff
		if err = sessionForUpdate.Run("python /root/sshd.py"); err != nil {
			logger.Printf("%s command execute output: %s\nError: %s\n", ip, outUpBuff.String(), errUpBuff.String())
			continue
		}
		fmt.Println(ip + " copy mission complete")
	}
	return nil
}

func delete() error {
	logger := createLogger()
	Hosts := readFile()
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		ip := items.Host
		szPassword := items.Password
		config := &ssh.ClientConfig{
			User: items.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(szPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		client, err := ssh.Dial("tcp", ip, config)
		if err != nil {
			logger.Printf("%s %s Maybe is the password error or is the host not avaliable\n", ip, err.Error())
			continue
		}
		defer client.Close()
		sessionForRM, sRMErr := client.NewSession()
		if sRMErr != nil {
			logger.Printf("%s Failed to create session: %s\n" + ip + sRMErr.Error())
			continue
		}
		defer sessionForRM.Close()
		var outRMBuff, errRMBuff bytes.Buffer
		sessionForRM.Stdout = &outRMBuff
		sessionForRM.Stderr = &errRMBuff
		if err = sessionForRM.Run("rm /root/sshd.py /root/sshd -f"); err != nil {
			logger.Printf("%s command execute output: %s\nError: %s\n", ip, outRMBuff.String(), errRMBuff.String())
			continue
		}
		fmt.Println(ip + " delete mission complete")
	}
	return nil
}

func scopy(remoteDir string) error {
	logger := createLogger()
	Hosts := readFile()
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		ip := items.Host
		szPassword := items.Password
		config := &ssh.ClientConfig{
			User: items.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(szPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		conn, err := ssh.Dial("tcp", ip, config)
		if err != nil {
			logger.Println(ip + err.Error())
			continue
		}
		defer conn.Close()
		sftpClient, sErr := sftp.NewClient(conn)
		if sErr != nil {
			logger.Println(ip + sErr.Error())
			continue
		}
		defer sftpClient.Close()
		sshdPyFile, fErr := os.Open("./sshd.py")
		if fErr != nil {
			logger.Println(ip + fErr.Error())
			continue
		}
		defer sshdPyFile.Close()
		sshdFile, sErr := os.Open("./sshd")
		if sErr != nil {
			logger.Println(ip + sErr.Error())
			continue
		}
		defer sshdFile.Close()
		// this is sshd.py file transfer.
		dstPYFile, pyErr := sftpClient.Create(remoteDir + "/sshd.py")
		if pyErr != nil {
			logger.Println(ip + pyErr.Error())
			continue
		}
		defer dstPYFile.Close()
		cpyData, cpyErr := io.Copy(dstPYFile, sshdPyFile)
		if cpyErr != nil {
			logger.Println(cpyErr.Error())
			continue
		}
		fmt.Printf("%s: %s %d has copies\n", ip, sshdPyFile.Name(), cpyData)
		// this is sshd file transfer.
		dstFile, pyErr := sftpClient.Create(remoteDir + "/sshd")
		if pyErr != nil {
			logger.Println(ip + pyErr.Error())
			continue
		}
		defer dstFile.Close()
		cpData, cpErr := io.Copy(dstFile, sshdFile)
		if cpErr != nil {
			logger.Println(ip + cpErr.Error())
			continue
		}
		fmt.Printf("%s: %s %d has copies\n", ip, sshdFile.Name(), cpData)
	}
	return nil
}

func main() {
	// test()
	execute()
}
