package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"service/api"
)

var (
	SingleCredentials string
	SingleParams      []string
	Verbose			  bool
)

var callCmd = &cobra.Command{
	Use:   "call [API NAME] [ENDPOINT]",
	Short: "Calls an api immediately",
	Args:  cobra.ExactArgs(2),
	Run:   callHandler,
}

func init() {
	addCredsFlag(callCmd, &SingleCredentials)
	addParamsFlag(callCmd, &SingleParams)
	addVerboseFlag(callCmd, &Verbose)
	RootCmd.AddCommand(callCmd)
}

func callHandler(cmd *cobra.Command, args []string) {
	config, err := api.GetApiConfig(args[0])
	if err != nil {
		log.Fatal(err)
	}
	arguments := api.GetArgumentsFromSlice(SingleParams)
	res, err := config.CallWGet(args[1], SingleCredentials, arguments[0])
	fmt.Println(res.Print(Verbose))
}
