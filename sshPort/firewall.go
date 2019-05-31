package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type (
	// Zone for firewall xml.
	Zone struct {
		// XMLName     xml.Name `xml:"-"`
		Text        string `xml:",chardata"`
		Short       string `xml:"short"`
		Description string `xml:"description"`
		Service     []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"service"`
		Port []struct {
			Text     string `xml:",chardata"`
			Protocol string `xml:"protocol,attr"`
			Port     string `xml:"port,attr"`
		} `xml:"port"`
		Rule []struct {
			Text   string `xml:",chardata"`
			Family string `xml:"family,attr"`
			Source struct {
				Text    string `xml:",chardata"`
				Address string `xml:"address,attr"`
			} `xml:"source"`
			Port struct {
				Text     string `xml:",chardata"`
				Protocol string `xml:"protocol,attr"`
				Port     string `xml:"port,attr"`
			} `xml:"port"`
			Accept string `xml:"accept"`
		} `xml:"rule"`
	}
)

func xmlReader() (*Zone, error) {
	file, err := os.Open("/etc/firewalld/zones/public.xml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	// buffer := xml.NewDecoder(bytes.NewReader(buff))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fZone := &Zone{}
	err = xml.Unmarshal(buff, &fZone)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return fZone, nil
}

func readConfig() (*Zone, error) {
	file, err := os.Open("./sshPort/rules.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, _ := ioutil.ReadAll(file)
	result := &Zone{}
	// fmt.Println(buffer)
	err = xml.Unmarshal(buffer, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getLogging() *log.Logger {
	file, err := os.Create("./firewalld.log")
	if err != nil {
		log.Fatalf("failed to create log file. %s\n", err.Error())
	}
	defer file.Close()
	return log.New(file, "", log.LstdFlags|log.Llongfile)
}

func main() {
	logger := getLogging()
	config, err := readConfig()
	if err != nil {
		logger.Println("config is not validate xml file.")
	}

	// TODO on linux debugging. uncommand below line.
	cmd := exec.Command("mv", "/etc/firewalld/zones/public.xml", "/etc/firewalld/zones/public.xml.bak")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() err is %s\n", stderr.String())
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)

	// file, _ := xml.MarshalIndent(parser, "", " ")
	// file = bytes.Replace(file, []byte("<accept></accept>"), []byte("<accept/>"), -1)
	// file = bytes.Replace(file, []byte("></port>"), []byte("/>"), -1)
	// file = bytes.Replace(file, []byte("></source>"), []byte("/>"), -1)

	buffer, _ := xml.MarshalIndent(config, "", " ")
	buffer = bytes.Replace(buffer, []byte("&#xA;  "), []byte(""), -1)
	buffer = bytes.Replace(buffer, []byte("&#xA;"), []byte(""), -1)
	buffer = bytes.Replace(buffer, []byte("&gt;"), []byte(""), -1)
	// TODO modify this to /etc/firewalld/zones/public.xml
	err = ioutil.WriteFile("/etc/firewalld/zones/public.xml", buffer, 0644)
	if err != nil {
		logger.Printf("write file error. %s", err.Error)
	}
}
