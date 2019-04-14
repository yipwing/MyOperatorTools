package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

const (
	usage = "backup /var/log/big /storage/backup 2022-01-01.log"
)

// TODO finish code.
func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	anwser, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("some error contact author")
	}
	if anwser == "y" || anwser == "yes" {
		self := os.Args[0]
		PathSeparator := "/"
		if runtime.GOOS == "windows" {
			PathSeparator = "\\"
		}
		program := strings.Split(self, PathSeparator)[len(strings.Split(self, PathSeparator))-1]

		checkProgam := "backup"
		if runtime.GOOS == "windows" {
			checkProgam = "backup.exe"
		}
		if program != checkProgam {
			fmt.Println("do not change program name")
			return
		}
		args := os.Args[1:]
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		moveFile(ctx, args[1], args[2], args[3])
	}
}

// TODO finish move file.
func moveFile(ctx context.Context, source, dest, except string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		break
	default:
		break
	}
}
