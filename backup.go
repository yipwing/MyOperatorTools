package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO finish code.
func main() {
	source := flag.String("source", "", "full path of directory")
	dest := flag.String("dest", "", "full path of destination")
	except := flag.String("except", "", "except file")
	flag.Parse()
	if len(*source) <= 0 || len(*dest) <= 0 || len(*except) <= 0 {
		fmt.Println("parameter empty")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("are u sure u want move these files?(y(yes)/any)")
	anwser, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimRight(anwser, "\r\n")) == "y" || strings.ToLower(strings.TrimRight(anwser, "\r\n")) == "yes" {
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
			fmt.Println("do not change the program name")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		moveFile(ctx, *source, *dest, *except)
	}
}

// TODO finish move file.
func moveFile(ctx context.Context, source, dest, except string) error {
	PathSeparator := "/"
	if runtime.GOOS == "windows" {
		PathSeparator = "\\"
	}

	select {
	case <-time.After(500 * time.Nanosecond):
		err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == source || info.Name() == except {
				return nil
			}
			if _, err = os.Stat(dest); os.IsNotExist(err) {
				cErr := os.MkdirAll(dest, 755)
				if cErr != nil {
					return cErr
				}
			}

			oErr := os.Rename(path, dest+PathSeparator+info.Name())
			if oErr != nil {
				return oErr
			}
			fmt.Println(path + "move to " + dest)
			return nil
		})
		if err != nil {
			return err
		}
	case <-ctx.Done():
		fmt.Println("halt moveFile")
	}
	return nil
}
