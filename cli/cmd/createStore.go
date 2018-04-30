package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thbishop/kv/cli/client"
)

var createStoreCmd = &cobra.Command{
	Use:   "create-store",
	Short: "Creates a new store",
	Long: `Creates a new store to store key/values. For example:

kv create-store --store-name my-store

`,
	Run: func(cmd *cobra.Command, args []string) {
		storeName := cmd.Flag("store-name").Value.String()
		err := client.CreateStore(storeName)
		if err != nil {
			fmt.Printf("Error creating store: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Store '%s' created successfully\n", storeName)
	},
}

func init() {
	rootCmd.AddCommand(createStoreCmd)
	createStoreCmd.Flags().StringP("store-name", "", "", "Name of the store to create")
	createStoreCmd.MarkFlagRequired("store-name")
}
