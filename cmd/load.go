package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var FromFlag string
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads api configuration files from a specified folder",
	Args:  cobra.ExactArgs(1),
	Run:   loadHandler,
}

func loadHandler(cmd *cobra.Command, args []string) {
	folder := fmt.Sprintf("%s/configs", rootFolder)
	// Check if the folder exists, if not, create one
	if _, err := os.Stat(folder); err != nil {
		if err := os.Mkdir(folder, 755); err != nil {
			color.Red("Cannot generate configurations folder (configs)")
			return
		}
	}
	if _, err := os.Stat(FromFlag); err != nil {
		panic(color.New(color.FgRed).Sprintf("Cant access configuration folder %s", FromFlag))
	}
	files, err := ioutil.ReadDir(FromFlag)
	if err != nil {
		panic(color.New(color.FgRed).Sprintf("Cant read files from configuration folder %s", FromFlag))
	}
	for _, file := range files {
		fileLocation := fmt.Sprintf("%s/%s", FromFlag, file.Name())
		contents, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			panic(color.New(color.FgRed).Sprintf("Cant read file %s from configuration folder %s", file.Name(), FromFlag))
		}
		if err := ioutil.WriteFile(fmt.Sprintf("%s/configs/%s", rootFolder, file.Name()), contents, 755); err != nil {
			panic(color.New(color.FgRed).Sprintf("Cant generate file %s", file.Name()))
		}
	}
}

func init() {
	loadCmd.Flags().StringVarP(
		&FromFlag,
		"from",
		"f",
		"",
		"Specify a folder to fetch the api configurations from")
	RootCmd.AddCommand(loadCmd)
}
