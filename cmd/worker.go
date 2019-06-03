package cmd

import (
	"fmt"

	"github.com/alfg/enc/api/worker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the worker.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting worker...")
		startWorkers()
	},
}

func startWorkers() {
	// Create Workers.
	worker.NewWorker(1, 100) // TODO: Use flags.
}
