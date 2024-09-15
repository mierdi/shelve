package cmd

import (
	"fmt"
	"os"
	"shelve/git"

	"github.com/spf13/cobra"
)

var cmdClear = &cobra.Command{
	Use: "clear",
	Run: clearWorkspace,
}

func init() {
	MainCommand.AddCommand(cmdClear)
}

func clearWorkspace(cmd *cobra.Command, args []string) {
	if err := git.ClearWorkspace(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
