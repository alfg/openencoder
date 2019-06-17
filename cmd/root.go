package cmd

import (
	"fmt"
	"os"

	"github.com/alfg/enc/api/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encoding API and worker.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "default", "Config YAML")

	config.LoadConfig(cfgFile)
	fmt.Println(config.Get())
}

// Execute starts cmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
