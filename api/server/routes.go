package server

import "github.com/gin-gonic/gin"

func registerRoutes(r *gin.Engine) {

	// JWT middleware.
	authMiddlware := jwtMiddleware()
	r.POST("/api/register", registerHandler)
	r.POST("/api/login", authMiddlware.LoginHandler)
	r.GET("/api/refresh-token", authMiddlware.RefreshHandler)
	r.GET("/api/", indexHandler)

	// API routes.
	api := r.Group("/api")
	api.Use(authMiddlware.MiddlewareFunc())
	{
		// User.
		api.GET("/user", getUserHandler)
		api.PUT("/user", updateUserHandler)

		// S3.
		api.GET("/s3/list", s3ListHandler)

		// Jobs.
		api.POST("/jobs", createJobHandler)
		api.GET("/jobs", getJobsHandler)
		api.GET("/jobs/:id", getJobsByIDHandler)
		api.PUT("/jobs/:id", updateJobByIDHandler)
		api.GET("/jobs/:id/status", getJobStatusByIDHandler)
		api.POST("/jobs/:id/cancel", cancelJobByIDHandler)
		api.POST("/jobs/:id/restart", restartJobByIDHandler)

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

		// Presets.
		api.POST("/presets", createPresetHandler)
		api.GET("/presets", getPresetsHandler)
		api.GET("/presets/:id", getPresetByIDHandler)
		api.PUT("/presets/:id", updatePresetByIDHandler)

		// Settings.
		api.GET("/settings", settingsHandler)
		api.PUT("/settings", updateSettingsHandler)
	}
}
