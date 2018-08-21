package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/theclocker/buckit/api"
	"log"
)

var (
	BulkParams      []string
	BulkCredentials string
)

var bulkCmd = &cobra.Command{
	Use:   "bulk [API NAME] [ENDPOINT]",
	Short: "Creates and runs a queue of api calls with the given flags",
	Args:  cobra.ExactArgs(2),
	Run:   bulkHandler,
}

func bulkHandler(cmd *cobra.Command, args []string) {
	config, err := api.GetApiConfig(args[0])
	if err != nil {
		log.Fatal(err)
	}
	arguments := api.GetArgumentsFromSlice(BulkParams)
	paramsArr := []map[string]string{
		arguments,
		arguments,
		arguments,
	}
	resChan := config.BulkCallWGet(args[1], BulkCredentials, paramsArr)
	for a := 0; a < len(paramsArr); a++ {
		res := <-resChan
		fmt.Println(res.Print())
	}
}

func init() {
	addCredsFlag(bulkCmd, &BulkCredentials)
	addParamsFlag(bulkCmd, &BulkParams)
	RootCmd.AddCommand(bulkCmd)
}
