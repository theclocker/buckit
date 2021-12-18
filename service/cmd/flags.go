package cmd

import "github.com/spf13/cobra"

func addCredsFlag(cmd *cobra.Command, credsVar *string) {
	cmd.Flags().StringVarP(
		credsVar,
		"credentials",
		"c",
		"default",
		"Which credentials should the call use")
}

// addParamsFlag defines the flag for parameters, can be used in bulk using the syntax:
// -p "key:value" -p "key:value"
func addParamsFlag(cmd *cobra.Command, paramsVar *[]string) {
	cmd.Flags().StringArrayVarP(
		paramsVar,
		"parameters",
		"p",
		[]string{},
		"Takes in a json object and uses it as api call parameters")
}

func addVerboseFlag(cmd *cobra.Command, verboseVar *bool) {
	cmd.Flags().BoolVarP(
		verboseVar,
		"verbose",
		"v",
		false,
		"Enter true false to see the response from the API, useful for debugging purposes.")
}