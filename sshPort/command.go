package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls")
	fmt.Println(cmd.Run())
}
