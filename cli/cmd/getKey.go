package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thbishop/kv/cli/client"
)

var getKeyCmd = &cobra.Command{
	Use:   "get-key",
	Short: "Gets the value of the key",
	Long: `Gets the value of the key. Example:

kv get-key --store-name my-store --key-name key1

`,
	Run: func(cmd *cobra.Command, args []string) {
		storeName := cmd.Flag("store-name").Value.String()
		keyName := cmd.Flag("key-name").Value.String()

		value, err := client.GetKey(storeName, keyName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting key: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", value)
	},
}

func init() {
	rootCmd.AddCommand(getKeyCmd)
	getKeyCmd.Flags().StringP("store-name", "", "", "Name of the store to create the key in")
	getKeyCmd.Flags().StringP("key-name", "", "", "Name of the key to set")

	getKeyCmd.MarkFlagRequired("store-name")
	getKeyCmd.MarkFlagRequired("key-name")
}
