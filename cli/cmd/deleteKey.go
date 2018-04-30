package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thbishop/kv/cli/client"
)

var deleteKeyCmd = &cobra.Command{
	Use:   "delete-key",
	Short: "Deletes a key",
	Long: `Deletes the desired key. For example:

kv delete-key --store-name my-store --key-name key1

** NOTE **: If the key does not exist, delete-key will *NOT* exit with an error. This is not considered an error as the intention is for the key to no longer exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		storeName := cmd.Flag("store-name").Value.String()
		keyName := cmd.Flag("key-name").Value.String()

		err := client.DeleteKey(storeName, keyName)
		if err != nil {
			if client.IsNotFoundError(err) {
				fmt.Printf("Key '%s' not found\n", keyName)
				os.Exit(0)
			}

			fmt.Printf("Error deleting key: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Key '%s' deleted successfully\n", keyName)
	},
}

func init() {
	rootCmd.AddCommand(deleteKeyCmd)
	deleteKeyCmd.Flags().StringP("store-name", "", "", "Name of the store to create the key in")
	deleteKeyCmd.Flags().StringP("key-name", "", "", "Name of the key to set")

	deleteKeyCmd.MarkFlagRequired("store-name")
	deleteKeyCmd.MarkFlagRequired("key-name")
}
