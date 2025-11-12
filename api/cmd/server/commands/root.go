package commands

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "beta-be",
	Short: "Run the beta-be command",
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "configs/.env", "Configuration file to use.")
}
