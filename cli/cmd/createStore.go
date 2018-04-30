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

kv create-store --name my-store

`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value.String()
		err := client.CreateStore(name)
		if err != nil {
			fmt.Printf("Error creating store: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Store '%s' created successfully\n", name)
	},
}

func init() {
	rootCmd.AddCommand(createStoreCmd)
	createStoreCmd.Flags().StringP("name", "", "", "Name of the store to create")
	createStoreCmd.MarkFlagRequired("name")
}
