package main

import (
	"fmt"
	"os"

	"github.com/hossein1376/spallet/cmd/spallet/command"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
