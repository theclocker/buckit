package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"service/api"
)

var (
	BulkParams      []string
	BulkCredentials string
	BulkVerbose			bool
)

var bulkCmd = &cobra.Command{
	Use:   "bulk [API NAME] [ENDPOINT]",
	Short: "Creates and runs a queue of api calls with the given flags",
	Args:  cobra.ExactArgs(2),
	Run:   bulkHandler,
}

func init() {
	addCredsFlag(bulkCmd, &BulkCredentials)
	addParamsFlag(bulkCmd, &BulkParams)
	addVerboseFlag(bulkCmd, &BulkVerbose)
	RootCmd.AddCommand(bulkCmd)
}

func bulkHandler(_ *cobra.Command, args []string) {
	config, err := api.GetApiConfig(args[0])
	if err != nil {
		log.Fatal(err)
	}
	arguments := api.GetArgumentsFromSlice(BulkParams)
	resChan := config.BulkCallWGet(args[1], BulkCredentials, arguments)
	for a := 0; a < len(arguments); a++ {
		res := <-resChan
		fmt.Println(res.Print(BulkVerbose))
	}
}
