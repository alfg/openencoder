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

func startWorkers() {
	// Create Workers.
	worker.NewWorker()
}

func main() {
	var (
		config = flag.String("config", "./default.yml", "Config YAML")
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

	startWorkers()
	startServer()
}
