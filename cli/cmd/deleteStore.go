package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thbishop/kv/cli/client"
)

var deleteStoreCmd = &cobra.Command{
	Use:   "delete-store",
	Short: "Deletes a store",
	Long: `Deletes the desired store. This will also delete *all* keys that exist in the store. For example:

kv delete-store --store-name my-store

** NOTE **: If the store does not exist, delete-store will *NOT* exit with an error. This is not considered an error as the intention is for the store to no longer exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		storeName := cmd.Flag("store-name").Value.String()

		err := client.DeleteStore(storeName)
		if err != nil {
			if client.IsNotFoundError(err) {
				fmt.Printf("Store '%s' not found\n", storeName)
				os.Exit(0)
			}

			fmt.Printf("Error deleting store: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Store '%s' deleted successfully\n", storeName)
	},
}

func init() {
	rootCmd.AddCommand(deleteStoreCmd)
	deleteStoreCmd.Flags().StringP("store-name", "", "", "Name of the store to create the store in")

	deleteStoreCmd.MarkFlagRequired("store-name")
}
