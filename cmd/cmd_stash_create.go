package cmd

import (
	"fmt"
	"os"
	"shelve/git"

	"github.com/spf13/cobra"
)

var cmdCreate = &cobra.Command{
	Use: "c",
	Run: createStash,
}

func init() {
	MainCommand.AddCommand(cmdCreate)
}

func createStash(cmd *cobra.Command, args []string) {
	var stashName string

	if len(args) > 0 {
		stashName = args[0]
	}

	if l, err := git.CreateStash(stashName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if len(l) > 0 {
		fmt.Println(l)
	}
}
