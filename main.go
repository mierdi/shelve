package main

import (
	"fmt"
	"os"
	"shelve/cmd"
)

func main() {
	cmd.MainCommand.SetArgs(os.Args[1:])

	if err := cmd.MainCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
