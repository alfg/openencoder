package server

import "github.com/gin-gonic/gin"

func registerRoutes(r *gin.Engine) {
	// JWT middleware.
	authMiddlware := setJWT()
	r.POST("/api/login", authMiddlware.LoginHandler)
	r.GET("/api/refresh-token", authMiddlware.RefreshHandler)

	// API routes.
	api := r.Group("/api")
	api.Use(authMiddlware.MiddlewareFunc())
	{
		// Index.
		api.GET("/", indexHandler)

		// S3.
		api.GET("/s3/list", s3ListHandler)

		// Profiles.
		api.GET("/profiles", profilesHandler)

		// Jobs.
		api.POST("/jobs", createJobHandler)
		api.GET("/jobs", getJobsHandler)
		api.GET("/jobs/:id", getJobsByIDHandler)
		api.PUT("/jobs/:id", updateJobByIDHandler)

		// Stats.
		api.GET("/stats", getStatsHandler)

		// Worker info.
		api.GET("/worker/queue", workerQueueHandler)
		api.GET("/worker/pools", workerPoolsHandler)
		api.GET("/worker/busy", workerBusyHandler)

		// Machines.
		api.GET("/machines", machinesHandler)
		api.POST("/machines", createMachineHandler)
		api.DELETE("/machines", deleteMachineByTagHandler)
		api.DELETE("/machines/:id", deleteMachineHandler)
		api.GET("/machines/regions", listMachineRegionsHandler)
		api.GET("/machines/sizes", listMachineSizesHandler)
	}

	// Auth.
	// api.POST("/register", registerHandler)
	// api.POST("/login", loginHandler)
}
