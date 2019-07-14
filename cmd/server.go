package cmd

import (
	"fmt"
	"runtime"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server...")
		startServer()
	},
}

func configRuntime() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Running with %d CPUs\n", numCPU)
}
func startServer() {

	// Server config.
	serverCfg := &server.Config{
		ServerPort:  config.Get().Port,
		RedisHost:   config.Get().RedisHost,
		RedisPort:   config.Get().RedisPort,
		Namespace:   config.Get().WorkerNamespace,
		JobName:     config.Get().WorkerJobName,
		Concurrency: config.Get().WorkerConcurrency,
	}

	// Create HTTP Server.
	configRuntime()
	server.NewServer(*serverCfg)
}
