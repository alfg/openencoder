package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var profilesFile string

var rootCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encoding API and worker.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "default", "Config YAML")
	rootCmd.PersistentFlags().StringVar(&profilesFile, "profiles", "profiles", "Profiles YAML")
}

// Execute starts cmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
