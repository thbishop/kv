package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thbishop/kv/cli/client"
)

var setKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Sets the key",
	Long: `Sets the key with the desired value. If the key is new, it will be created. If it is an existing key, the value will be overwritten. For example:

kv set-key --store-name my-store --key-name key1 --key-value foo

`,
	Run: func(cmd *cobra.Command, args []string) {
		storeName := cmd.Flag("store-name").Value.String()
		keyName := cmd.Flag("key-name").Value.String()
		keyValue := cmd.Flag("key-value").Value.String()
		err := client.SetKey(storeName, keyName, keyValue)
		if err != nil {
			fmt.Printf("Error setting key: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Key '%s' set successfully\n", keyName)
	},
}

func init() {
	rootCmd.AddCommand(setKeyCmd)
	setKeyCmd.Flags().StringP("store-name", "", "", "Name of the store to create the key in")
	setKeyCmd.Flags().StringP("key-name", "", "", "Name of the key to set")
	setKeyCmd.Flags().StringP("key-value", "", "", "Value of the key")

	setKeyCmd.MarkFlagRequired("store-name")
	setKeyCmd.MarkFlagRequired("key-name")
	setKeyCmd.MarkFlagRequired("key-value")
}
