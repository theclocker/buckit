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
)

var callCmd = &cobra.Command{
	Use:   "call [API NAME] [ENDPOINT]",
	Short: "Calls an api immediately",
	Args:  cobra.ExactArgs(2),
	Run:   callHandler,
}

func callHandler(cmd *cobra.Command, args []string) {
	config, err := api.GetApiConfig(args[0])
	if err != nil {
		log.Fatal(err)
	}
	arguments := api.GetArgumentsFromSlice(SingleParams)
	res, err := config.CallWGet(args[1], SingleCredentials, arguments)
	fmt.Println(res.Print())
}

func addCredsFlag(cmd *cobra.Command, credsVar *string) {
	cmd.Flags().StringVarP(
		credsVar,
		"credentials",
		"c",
		"default",
		"Which credentials should the call use")
}

func addParamsFlag(cmd *cobra.Command, paramsVar *[]string) {
	cmd.Flags().StringArrayVarP(
		paramsVar,
		"parameters",
		"p",
		[]string{},
		"Takes in a json object and uses it as api call parameters")
}

func init() {
	addCredsFlag(callCmd, &SingleCredentials)
	addParamsFlag(callCmd, &SingleParams)
	RootCmd.AddCommand(callCmd)
}
