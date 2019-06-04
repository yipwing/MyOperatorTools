package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	ownHost struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
)

func readFile() []ownHost {
	file, err := os.Open("./openssh/hosts.json")
	if err != nil {
		fmt.Println("info:", "failed to open json file")
	}
	defer file.Close()
	result := []ownHost{}
	buffer, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		fmt.Println("info:", "cannot read hosts.json file or is not validate json format.")
	}
	fmt.Println("has:", len(result))
	return result
}

func main() {
	fmt.Println(readFile())
}
