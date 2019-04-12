package main

import (
	"context"
	"fmt"
	"os"
)

// TODO finish code.
func main() {
	self := os.Args[0]
	fmt.Println(self)
	if self != "backup" {
		fmt.Println("change the program name?")
		return
	}
	args := os.Args[1:]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go moveFile(ctx, args[1], args[2], args[3])
}

func moveFile(ctx context.Context, source, dest, except string) {

}
