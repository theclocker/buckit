package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theclocker/buckit/api"
	"fmt"
)

var callCmd = &cobra.Command{
	Use: "call [API NAME] [ENDPOINT]",
	Short: "Call an api immediately",
	Args: cobra.ExactArgs(2),
	Run: callHandler,
}

func callHandler(cmd *cobra.Command, args []string) {
	config, _ := api.GetApiConfig("openweathermap")
	fmt.Printf("%v\n", config)
	//for _, arg := range args {
	//	fmt.Println(arg)
	//}
}

func init() {
	RootCmd.AddCommand(callCmd)
}
