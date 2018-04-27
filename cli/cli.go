package main

import(
	"fmt"
	"os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "kv",
  Short: "kv is the cli for the kv service",
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
  Version: "0.1.0",
}

func executeCli() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
