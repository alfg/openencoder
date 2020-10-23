package server

import (
	"github.com/alfg/openencoder/api/config"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
)

func workerQueueHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

	queues, err := client.Queues()
	if err != nil {
		log.Error(err)
	}
	c.JSON(200, queues)
}

func workerPoolsHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

	resp, err := client.WorkerPoolHeartbeats()
	if err != nil {
		log.Error(err)
	}
	c.JSON(200, resp)
}

func workerBusyHandler(c *gin.Context) {
	client := work.NewClient(config.Get().WorkerNamespace, redisPool)

	observations, err := client.WorkerObservations()
	if err != nil {
		log.Error(err)
	}

	var busyObservations []*work.WorkerObservation
	for _, ob := range observations {
		if ob.IsBusy {
			busyObservations = append(busyObservations, ob)
		}
	}
	c.JSON(200, busyObservations)
}
