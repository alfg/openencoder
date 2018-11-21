package main

import (
	_ "expvar"
	"flag"
	"fmt"
	_ "net/http/pprof"

	"github.com/alfg/enc/helpers"
	"github.com/alfg/enc/server"
	"github.com/alfg/enc/worker"
)

func main() {
	var (
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		port         = flag.String("port", "8080", "The server port")
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

	// Get workflow configs.
	helpers.LoadConfig(*config)
	// helpers.LoadConfig()
	fmt.Println(helpers.C)

	// Create Workers.
	worker.NewWorker(*maxWorkers, *maxQueueSize)

	// Create HTTP Server.
	server.NewServer(*port)

}
