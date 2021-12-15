package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var rootFolder = "/.buckit"
var RootCmd = &cobra.Command{
	Use:   "buckit",
	Short: "Manage API's without actually managing anything",
	Long:  `Manage API's without actually managing anything`,
}

func Execute() {
	generateRootFolder(rootFolder)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateRootFolder(folder string) {
	_, err := os.Stat(folder)
	if err != nil {
		color.Yellow("Folder %s does not exist, generating folder\n", folder)
	} else {
		return
	}
	if err := os.Mkdir(folder, 755); err != nil {
		color.Red("Can't generate root folder required for logging and maintain software output")
	}
}

func generateSupportingFolder(folder string) {
	if err := os.Mkdir(fmt.Sprintf("%s/logs", folder), 755); err != nil {
		color.Red("Cannot generate logs folder (logs)")
	}
}
