package main

import (
	_ "expvar"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"runtime"

	"github.com/alfg/enc/config"
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
	port := config.GetPort()
	server.NewServer(port)
}

func startWorkers() {
	// Create Workers.
	worker.NewWorker()
}

func main() {
	var (
		configFile = flag.String("configFile", "./config/default.yml", "Config YAML")
	)
	flag.Parse()

	config.LoadConfig(*configFile)

	configRuntime()
	startWorkers()
}
