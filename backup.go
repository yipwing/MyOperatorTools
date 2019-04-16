package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	usage = "Usage: \nbackup /var/log/big /storage/backup 2022-01-01.log\n parameter 1 is a source directory\n parameter 2 is destination of directory\n parameter 2 is except file."
)

// TODO finish code.
func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("are u sure u want move these files?(yes/no)")
	anwser, _ := reader.ReadString('\n')
	fmt.Print(anwser)
	// anwser = strings.Replace(anwser, "\r\n", "", -1)
	fmt.Println(os.Args)
	if strings.TrimRight(anwser, "\r\n") == "y" || strings.TrimRight(anwser, "\r\n") == "yes" {
		args := os.Args[1:]
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		fmt.Println(moveFile(ctx, args[0], args[1], args[2]))
	}
}

// TODO finish move file.
func moveFile(ctx context.Context, source, dest, except string) error {
	// files, err := ioutil.ReadDir(source)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, info.Size())
		return nil
	})
	if err != nil {
		return err
	}
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// 	if f.IsDir() {
	// 		for _, Recursive := range
	// 	}
	// }
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return nil
	}
}
