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
	defer file.Close()
	return log.New(file, "", log.LstdFlags|log.Llongfile)
}

func readFile() []ownHost {
	file, err := os.Open("./hosts.json")
	if err != nil {
		log.Println("failed to open json file")
		return nil
	}
	defer file.Close()
	result := []ownHost{}
	buffer, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		log.Println("cannot read hosts.json file or is not validate json format.")
		return nil
	}
	return result
}

func execute() error {
	logger := createLogger()
	Hosts := readFile()
	remoteDir := flag.String("remoteDir", "/root", "remote directory")
	flag.Parse()
	// fmt.Printf("%d mission ready to go \n", len(Hosts))
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		if dErr := items.delete(); dErr != nil {
			logger.Fatalf(dErr.Error())
			return dErr
		}
		err := items.scopy(*remoteDir)
		if err != nil {
			logger.Fatalf(err.Error())
			return err
		}
		szPassword := items.Password
		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password(szPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		client, err := ssh.Dial("tcp", items.Host, config)
		if err != nil {
			logger.Printf("%s %s Maybe is the password error or is the host not avaliable\n", items.Host, err.Error())
			continue
		}
		defer client.Close()
		sessionForUpdate, sUPErr := client.NewSession()
		if sUPErr != nil {
			logger.Printf("%s Failed to create session: %s\n" + items.Host + sUPErr.Error())
			continue
		}
		defer sessionForUpdate.Close()
		var outUpBuff, errUpBuff bytes.Buffer
		sessionForUpdate.Stdout = &outUpBuff
		sessionForUpdate.Stderr = &errUpBuff
		if err = sessionForUpdate.Run("chmod +x /root/firewall && /root/firewall"); err != nil {
			logger.Printf("%s command execute output: %s\nError: %s\n", items.Host, outUpBuff.String(), errUpBuff.String())
			continue
		}
		fmt.Println(items.Host + " copy mission complete")
	}
	return nil
}

func (item *ownHost) delete() error {
	logger := createLogger()
	config := &ssh.ClientConfig{
		User: item.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(item.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	client, err := ssh.Dial("tcp", item.Host, config)
	if err != nil {
		logger.Printf("%s %s Maybe is the password error or is the host not avaliable\n", item.Host, err.Error())
	}
	defer client.Close()
	sessionForRM, sRMErr := client.NewSession()
	if sRMErr != nil {
		logger.Printf("%s Failed to create session: %s\n" + item.Host + sRMErr.Error())
	}
	defer sessionForRM.Close()
	var outRMBuff, errRMBuff bytes.Buffer
	sessionForRM.Stdout = &outRMBuff
	sessionForRM.Stderr = &errRMBuff
	if err = sessionForRM.Run("rm /root/firewall -f"); err != nil {
		logger.Printf("%s command execute output: %s\nError: %s\n", item.Host, outRMBuff.String(), errRMBuff.String())
	}
	fmt.Println(item.Host + " delete mission complete")
	return nil
}

func (item *ownHost) scopy(remoteDir string) error {
	logger := createLogger()
	config := &ssh.ClientConfig{
		User: item.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(item.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	conn, err := ssh.Dial("tcp", item.Host, config)
	if err != nil {
		logger.Println(item.Host + err.Error())
	}
	defer conn.Close()
	sftpClient, sErr := sftp.NewClient(conn)
	if sErr != nil {
		logger.Println(item.Host + sErr.Error())
	}
	defer sftpClient.Close()
	firewall, fErr := os.Open("./firewall")
	if fErr != nil {
		logger.Println(item.Host + fErr.Error())
	}
	defer firewall.Close()
	// this is sshd.py file transfer.
	dstFile, pyErr := sftpClient.Create(remoteDir + "/firewall")
	if pyErr != nil {
		logger.Println(item.Host + pyErr.Error())
	}
	defer dstFile.Close()
	cpyData, cpyErr := io.Copy(dstFile, firewall)
	if cpyErr != nil {
		logger.Println(cpyErr.Error())
	}
	fmt.Printf("%s: %s %d has copies\n", item.Host, firewall.Name(), cpyData)
	return nil
}

func main() {
	// test()
	execute()
}
