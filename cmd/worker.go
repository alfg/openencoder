package cmd

import (
	"fmt"

	"github.com/alfg/enc/api/config"
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
	// config.LoadConfig(cfgFile)
	// fmt.Println(cfgFile)
	// fmt.Println("s3 region: ", config.Get().S3Region)

	// Worker config.
	workerCfg := &worker.Config{
		Host:        config.Get().RedisHost,
		Port:        config.Get().RedisPort,
		Namespace:   config.Get().WorkerNamespace,
		JobName:     config.Get().WorkerJobName,
		Concurrency: config.Get().WorkerConcurrency,
	}

	fmt.Println(workerCfg)

	// Create Workers.
	worker.NewWorker(*workerCfg)
}
