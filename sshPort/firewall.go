package main

import (
	"bytes"
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	RuleConfig struct {
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
	}
)

func xmlReader() (*Zone, error) {
	file, err := os.Open("E:\\codes\\MyOperatorTools\\sshPort\\debug.xml")
	// cmd := exec.Command("/bin/sh", "cp /etc/firewalld/zones/public.xml /etc/firewalld/zones/public.xml.bak")
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

func readConfig() ([]RuleConfig, error) {
	file, err := os.Open("./rules.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer,_ := ioutil.ReadAll(file)
	result := []RuleConfig{}
	err = xml.Unmarshal(buffer, &result)
	if err != nil {
		return nil, err
	}
	return result
}

func main() {
	parser, err := xmlReader()
	if err != nil {
		log.Println("parse xml error")
	}
	// if len(parser.Port) > 0 {
	// 	for i := 0; i < len(parser.Port); i++ {
	// 		parser.Port[i].Text = ""
	// 		parser.Port[i].Port = ""
	// 		parser.Port[i].Protocol = ""
	// 	}
	// }

	// TODO add rule to struct.
	// praser.Rule = append(parser.Rule)
	fmt.Println(parser)
	for i:=0;i<
	file, _ := xml.MarshalIndent(parser, "", " ")
	file = bytes.Replace(file, []byte("&#xA;  "), []byte(""), -1)
	file = bytes.Replace(file, []byte("&#xA;"), []byte(""), -1)
	file = bytes.Replace(file, []byte("<accept></accept>"), []byte("<accept/>"), -1)
	file = bytes.Replace(file, []byte("></port>"), []byte("/>"), -1)
	file = bytes.Replace(file, []byte("></source>"), []byte("/>"), -1)
	_ = ioutil.WriteFile("./sshPort/debug1.xml", file, 0644)
}
