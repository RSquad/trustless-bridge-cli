package cmd

import (
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message to a contract",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.PersistentFlags().StringP("address", "a", "", "Address of the contract")
	sendCmd.MarkFlagRequired("address")
}
