package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

// for host.
var (
	Hosts = []string{
		"172.17.101.129:3722|fanjie",
		"172.17.101.128:3722|fanjie",
		"172.17.101.127:3722|fanjie",
		"172.17.101.126:3722|fanjie",
		"172.17.101.125:3722|fanjie",
		"172.17.101.124:3722|fanjie",
		"172.17.101.123:3722|fanjie",
		"172.17.101.116:3722|fanjie",
		"172.17.101.115:3722|fanjie",
	}
)

// type (
// 	connData struct {
// 		auth         []ssh.AuthMethod
// 		addr         string
// 		clientConfig *ssh.ClientConfig
// 		client       *ssh.Client
// 		config       ssh.Config
// 		session      *ssh.Session
// 		err          error
// 	}
// )

func createLogger() *log.Logger {
	file, err := os.Create("./exec.log")
	if err != nil {
		log.Fatalln("failed to create log file.")
	}
	return log.New(file, "", log.LstdFlags|log.Llongfile)
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

// TODO test remote dir.  test different directory.
// func test() {
// 	config := &ssh.ClientConfig{
// 		User: "root",
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(szPassword),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}
// 	client, err := ssh.Dial("tcp", ipAddrs[0], config)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		logger.Fatalln("hanshake failed: unable to authenticate")
// 	}
// 	if err != nil {
// 		logger.Fatalln("Failed to dial: " + err.Error())
// 		return
// 	}
// 	session, sErr := client.NewSession()
// 	if sErr != nil {
// 		logger.Fatalln(ipAddrs[0] + " failed to create session: " + sErr.Error())
// 	}
// 	defer session.Close()
// 	var outBuff, errBuff bytes.Buffer
// 	session.Stdout = &outBuff
// 	session.Stderr = &errBuff
// 	if err = session.Run("ls -al /var/log"); err != nil {
// 		logger.Fatalln(ipAddrs[0] + "failed to run: " + err.Error())
// 	}
// 	logger.Fatalln(outBuff.String())
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

	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		item := strings.Split(items, "|")
		ip := item[0]
		szPassword := item[1]
		config := &ssh.ClientConfig{
			User: "root",
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
		logger.Println(ip + " copy mission complete")
	}
	return nil
}

func delete() error {
	logger := createLogger()
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		item := strings.Split(items, "|")
		ip := item[0]
		szPassword := item[1]
		config := &ssh.ClientConfig{
			User: "root",
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
		logger.Println(ip + " delete mission complete")
	}
	return nil
}

func scopy(remoteDir string) error {
	logger := createLogger()
	for _, items := range Hosts {
		time.Sleep(300 * time.Microsecond)
		item := strings.Split(items, "|")
		ip := item[0]
		szPassword := item[1]
		config := &ssh.ClientConfig{
			User: "root",
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

// ipAddrs = []string{
// 	"172.17.103.226:3722",
// 	"172.17.103.225:3722",
// 	"172.17.103.224:3722",
// 	"172.17.103.223:3722",
// 	"172.17.103.222:3722",
// 	"172.17.105.30:3722",
// 	"172.17.105.31:3722",
// 	"172.17.105.32:3722",
// 	"172.17.113.42:3722",
// 	"172.17.113.41:3722",
// 	"172.17.113.44:3722",
// 	"172.17.113.43:3722",
// 	"172.17.107.44:3722",
// 	"172.17.113.38:3722",
// 	"172.17.113.37:3722",
// 	"172.17.113.36:3722",
// 	"172.17.113.35:3722",
// 	"172.17.113.34:3722",
// 	"172.17.113.33:3722",
// 	"172.17.113.32:3722",
// 	"172.17.113.31:3722",
// 	"172.18.1.222:3722",
// 	"172.17.107.43:3722",
// 	"172.17.107.42:3722",
// 	"172.17.107.41:3722",
// 	"172.17.107.40:3722",
// 	"172.17.107.39:3722",
// 	"172.17.107.38:3722",
// 	"172.17.107.37:3722",
// 	"172.17.107.36:3722",
// 	"172.17.107.35:3722",
// 	"172.17.107.34:3722",
// 	"172.17.113.30:3722",
// 	"172.17.113.29:3722",
// 	"172.17.113.27:3722",
// 	"172.17.113.28:3722",
// 	"172.17.113.25:3722",
// 	"172.17.113.24:3722",
// 	"172.17.113.22:3722",
// 	"172.17.113.21:3722",
// 	"172.17.107.33:3722",
// 	"172.17.107.32:3722",
// 	"172.17.107.31:3722",
// 	"172.17.107.30:3722",
// 	"172.17.107.29:3722",
// 	"172.17.107.27:3722",
// 	"172.17.107.26:3722",
// 	"172.17.107.25:3722",
// 	"172.17.107.24:3722",
// 	"172.17.105.23:3722",
// 	"172.17.105.25:3722",
// 	"172.17.103.227:3722",
// 	"172.17.103.228:3722",
// 	"172.17.103.229:3722",
// 	"172.17.103.231:3722",
// 	"172.17.103.232:3722",
// 	"172.17.103.233:3722",
// 	"172.17.103.234:3722",
// 	"172.17.103.235:3722",
// 	"172.17.101.129:3722",
// 	"172.17.101.128:3722",
// 	"172.17.103.230:3722",
// 	"172.17.103.237:3722",
// 	"172.17.101.127:3722",
// 	"172.17.101.126:3722",
// 	"172.17.101.125:3722",
// 	"172.17.101.124:3722",
// 	"172.17.101.123:3722",
// 	"172.17.103.236:3722",
// 	"172.17.103.238:3722",
// 	"172.17.103.239:3722",
// 	"172.17.100.117:3722",
// 	"172.17.100.116:3722",
// 	"172.17.101.116:3722",
// 	"172.17.101.115:3722",
// 	"172.17.100.115:3722",
// 	"172.17.100.114:3722",
// 	"172.17.103.240:3722",
// 	"172.17.103.241:3722",
// 	"172.17.103.242:3722",
// 	"172.17.103.243:3722",
// 	"172.17.107.23:3722",
// 	"172.17.107.22:3722",
// 	"172.17.105.26:3722",
// 	"172.17.107.21:3722",
// 	"172.17.107.20:3722",
// 	"172.17.105.27:3722",
// 	"172.17.100.93:3722",
// 	"172.17.101.87:3722",
// 	"172.17.107.19:3722",
// 	"172.17.103.245:3722",
// 	"172.17.107.13:3722",
// 	"172.17.107.12:3722",
// 	"172.17.107.11:3722",
// 	"172.17.107.10:3722",
// 	"172.17.107.15:3722",
// 	"172.17.107.14:3722",
// 	"172.17.107.8:3722",
// 	"172.17.107.7:3722",
// 	"172.17.107.17:3722",
// 	"172.17.107.18:3722",
// 	"172.17.107.4:3722",
// 	"172.17.107.3:3722",
// 	"172.17.103.246:3722",
// 	"172.17.103.247:3722",
// 	"172.17.103.248:3722",
// 	"172.17.105.20:3722",
// 	"172.17.103.249:3722",
// 	"172.17.105.18:3722",
// 	"172.17.105.17:3722",
// 	"172.17.105.13:3722",
// 	"172.17.101.122:3722",
// 	"172.17.101.121:3722",
// 	"172.17.101.120:3722",
// 	"172.17.101.112:3722",
// 	"172.17.101.110:3722",
// 	"172.17.101.111:3722",
// 	"172.17.101.109:3722",
// 	"172.17.113.20:3722",
// 	"172.17.113.19:3722",
// 	"172.17.113.18:3722",
// 	"172.17.113.17:3722",
// 	"172.17.113.16:3722",
// 	"172.17.113.14:3722",
// 	"172.17.113.15:3722",
// 	"172.17.113.13:3722",
// 	"172.17.113.12:3722",
// 	"172.17.113.11:3722",
// 	"172.17.113.10:3722",
// 	"172.17.113.9:3722",
// 	"172.17.113.8:3722",
// 	"172.17.113.7:3722",
// 	"172.17.113.6:3722",
// 	"172.17.113.5:3722",
// 	"172.17.100.108:3722",
// 	"172.17.101.105:3722",
// 	"172.17.101.104:3722",
// 	"172.17.100.103:3722",
// 	"172.17.102.30:3722",
// 	"172.17.100.105:3722",
// 	"172.17.100.106:3722",
// 	"172.17.100.104:3722",
// 	"172.17.100.101:3722",
// 	"172.17.100.102:3722",
// 	"172.18.1.220:3722",
// 	"172.17.100.113:3722",
// 	"172.17.100.112:3722",
// 	"172.17.101.92:3722",
// 	"172.17.100.94:3722",
// 	"172.17.102.13:3722",
// 	"172.17.102.12:3722",
// 	"172.17.102.11:3722",
// 	"172.17.102.9:3722",
// 	"172.17.102.48:3722",
// 	"172.17.101.48:3722",
// 	"172.17.100.48:3722",
// 	"172.17.102.46:3722",
// 	"172.17.101.86:3722",
// 	"172.17.100.49:3722",
// 	"172.17.102.49:3722",
// 	"172.17.102.41:3722",
// 	"172.17.102.40:3722",
// 	"172.17.101.5:3722",
// 	"172.17.101.233:3722",
// 	"172.17.100.89:3722",
// 	"172.17.100.99:3722",
// 	"172.17.100.98:3722",
// 	"172.17.100.96:3722",
// 	"172.17.100.95:3722",
// 	"172.17.100.94:3722",
// 	"172.17.100.91:3722",
// 	"172.17.100.90:3722",
// 	"172.17.100.88:3722",
// 	"172.17.100.87:3722",
// 	"172.17.101.85:3722",
// 	"172.17.101.84:3722",
// 	"172.17.101.83:3722",
// 	"172.17.101.82:3722",
// 	"172.17.101.81:3722",
// 	"172.17.101.80:3722",
// 	"172.17.101.79:3722",
// 	"172.17.101.78:3722",
// 	"172.17.101.77:3722",
// 	"172.17.101.76:3722",
// 	"172.17.101.75:3722",
// 	"172.17.101.74:3722",
// 	"172.17.101.73:3722",
// 	"172.17.101.72:3722",
// 	"172.17.101.71:3722",
// 	"172.17.101.70:3722",
// 	"172.17.101.69:3722",
// 	"172.17.101.68:3722",
// 	"172.17.101.67:3722",
// 	"172.17.101.66:3722",
// 	"172.17.101.64:3722",
// 	"172.17.101.63:3722",
// 	"172.17.101.62:3722",
// 	"172.17.101.61:3722",
// 	"172.17.101.60:3722",
// 	"172.17.101.59:3722",
// 	"172.17.101.58:3722",
// 	"172.17.101.57:3722",
// 	"172.17.101.55:3722",
// 	"172.17.101.54:3722",
// 	"172.17.100.111:3722",
// 	"172.17.100.110:3722",
// 	"172.17.101.51:3722",
// 	"172.17.101.50:3722",
// 	"172.17.101.49:3722",
// 	"172.17.102.29:3722",
// 	"172.17.102.8:3722",
// 	"172.17.102.7:3722",
// 	"172.17.102.6:3722",
// 	"172.17.102.5:3722",
// 	"172.17.102.4:3722",
// 	"172.17.102.3:3722",
// 	"172.17.102.2:3722",
// 	"172.17.105.29:3722",
// }
