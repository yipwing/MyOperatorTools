package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type (
	// Zone for firewall xml.
	Zone struct {
		XMLName     xml.Name `xml:"zone"`
		Text        string   `xml:",chardata"`
		Short       string   `xml:"short"`
		Description string   `xml:"description"`
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
	file, err := os.Open("E:\\codes\\MyOperatorTools\\sshPort\\debug.xml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)
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

func main() {
	parser, err := xmlReader()
	if err != nil {
		log.Println("parse xml error")
	}
	// process xml
}
