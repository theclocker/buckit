package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets the whole application, deleting cache, logs and root folder",
	Run:   resetHandler,
}

func resetHandler(cmd *cobra.Command, args []string) {
	_, err := os.Stat(rootFolder)
	if err != nil {
		color.Red("Cant find root folder \"%s\", command exiting", rootFolder)
		return
	}
	if err := os.RemoveAll(rootFolder); err != nil {
		color.Red("Could not delete root folder \"%s\", command exiting", rootFolder)
	}
}

func init() {
	RootCmd.AddCommand(resetCmd)
}
