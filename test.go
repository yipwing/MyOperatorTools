package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type (
	ownHost struct {
		Host     string `json:"host"`
		Password string `json:"password"`
	}
)

func readFile() []ownHost {
	file, err := os.Open("./hosts.json")
	if err != nil {
		log.Println("failed to open json file")
	}
	defer file.Close()
	result := []ownHost{}
	buffer, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		log.Println("cannot read hosts.json file or is not validate json format.")
	}
	return result
}

func main() {
	fmt.Println(readFile())
}
