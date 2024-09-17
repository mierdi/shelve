package cmd

import (
	"fmt"
	"os"
	"shelve/git"

	"github.com/spf13/cobra"
)

var MainCommand = &cobra.Command{
	Run: start,
}

func start(cmd *cobra.Command, args []string) {
	if l, err := git.CreateStash(""); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if len(l) > 0 {
		fmt.Println(l)
	}
}
