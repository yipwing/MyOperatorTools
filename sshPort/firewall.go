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
		// Service     []struct {
		// 	Text string `xml:",chardata"`
		// 	Name string `xml:"name,attr"`
		// } `xml:"-"`
		// Port []struct {
		// 	Text     string `xml:",chardata"`
		// 	Protocol string `xml:"protocol,attr"`
		// 	Port     string `xml:"port,attr"`
		// } `xml:"-"`
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
	// Rules for all rule.
	Rules struct {
		Service []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"service"`
		Port []struct {
			Text     string `xml:",chardata"`
			Protocol string `xml:"protocol,attr"`
			Port     string `xml:"port,attr"`
		} `xml:"port"`
		RuleConfig []struct {
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
	// file, err := os.Open("E:\\codes\\MyOperatorTools\\sshPort\\debug.xml")
	file, err := os.Open("./sshPort/debug.xml")
	cmd := exec.Command("cp", "/etc/firewalld/zones/public.xml", "/etc/firewalld/zones/public.xml.bak")
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

func readConfig() (*Rules, error) {
	file, err := os.Open("./sshPort/rules.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, _ := ioutil.ReadAll(file)
	result := &Rules{}
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
	parser, err := xmlReader()
	if err != nil {
		logger.Println("parse xml error")
	}
	config, cErr := readConfig()
	if cErr != nil {
		logger.Println("config is not validate xml file.")
	}
	parser.Rule = config.RuleConfig
	// if len(parser.Port) > 0 {
	// 	for i := 0; i < len(parser.Port); i++ {
	// 		parser.Port[i].Text = ""
	// 		parser.Port[i].Port = ""
	// 		parser.Port[i].Protocol = ""
	// 	}
	// }

	// TODO add rule to struct.
	// praser.Rule = append(parser.Rule)
	// fmt.Println(parser)

	// file, _ := xml.MarshalIndent(parser, "", " ")
	// file = bytes.Replace(file, []byte("&#xA;  "), []byte(""), -1)
	// file = bytes.Replace(file, []byte("&#xA;"), []byte(""), -1)
	// file = bytes.Replace(file, []byte("<accept></accept>"), []byte("<accept/>"), -1)
	// file = bytes.Replace(file, []byte("></port>"), []byte("/>"), -1)
	// file = bytes.Replace(file, []byte("></source>"), []byte("/>"), -1)
	// _ = ioutil.WriteFile("./sshPort/debug1.xml", file, 0644)

	buffer, _ := xml.MarshalIndent(parser, "", " ")
	buffer = bytes.Replace(buffer, []byte("&#xA;  "), []byte(""), -1)
	buffer = bytes.Replace(buffer, []byte("&#xA;"), []byte(""), -1)
	err = ioutil.WriteFile("./sshPort/debug1.xml", buffer, 0644)
	if err != nil {
		logger.Printf("write file error. %s", err.Error)
	}
}
