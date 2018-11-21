package main

import (
	_ "expvar"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"runtime"

	"github.com/alfg/enc/helpers"
	"github.com/alfg/enc/server"
	"github.com/alfg/enc/worker"
)

func configRuntime() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Running with %d CPUs\n", numCPU)
}

func startServer() {
	// Get workflow configs.
	// // helpers.LoadConfig()
	// fmt.Println(helpers.C)

	// Create HTTP Server.
	port := helpers.GetPort()
	server.NewServer(port)
}

func startWorkers(maxWorkers *int, maxQueueSize *int) {
	// Create Workers.
	worker.NewWorker(*maxWorkers, *maxQueueSize)
}

func main() {
	var (
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		config       = flag.String("config", "./default.yml", "Config YAML")
	)
	flag.Parse()

	// Load config.
	// viper.SetConfigName("default")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(fmt.Errorf("fatal error config file: %s", err))
	// }

	// fmt.Println(viper.Get("settings.s3_bucket"))

	helpers.LoadConfig(*config)
	configRuntime()

	startWorkers(maxWorkers, maxQueueSize)
	startServer()
}
