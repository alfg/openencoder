package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encoding API and worker.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("stuff")
	},
}

// Execute starts cmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
