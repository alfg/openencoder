package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/machine"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
)

type machineRequest struct {
	Provider string `json:"provider" binding:"required"`
	Size     string `json:"size" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1,max=10"` // Max of 10 machines.
}

func machinesHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	setting, err := d.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(setting.Value)
	ctx := context.TODO()

	// Get list of machines from DO client.
	machines, err := client.ListDropletByTag(ctx, WorkerTag)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"machines": machines,
	})
}

func createMachineHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json machineRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := data.New()
	settings := db.Settings.GetSettings()

	token := types.GetSetting(types.DigitalOceanAccessToken, settings)
	region := types.GetSetting(types.DigitalOceanRegion, settings)
	vpc := types.GetSetting(types.DigitalOceanVPC, settings)

	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Create machine.
	machine, err := client.CreateDroplets(ctx, region, json.Size, vpc, json.Count)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	// TODO: Add resource to project?

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	d := data.New()
	token, err := d.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(token.Value)
	ctx := context.TODO()

	// Create machine.
	machine, err := client.DeleteDropletByID(ctx, id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineByTagHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	token, err := d.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(token.Value)
	ctx := context.TODO()

	// Create machine.
	err = client.DeleteDropletByTag(ctx, WorkerTag)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"deleted": true,
	})
	return
}

func listMachineRegionsHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	db := data.New()
	settings := db.Settings.GetSettings()

	token := types.GetSetting(types.DigitalOceanAccessToken, settings)
	region := types.GetSetting(types.DigitalOceanRegion, settings)
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Get list of machine regions from DO client.
	regions, err := client.ListRegions(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	// Filter regions by configured region from settings.
	var filteredRegions []machine.Region
	for _, r := range regions {
		if r.Slug == region {
			filteredRegions = append(filteredRegions, r)
		}
	}

	c.JSON(200, gin.H{
		"regions": filteredRegions,
	})
}

func listMachineSizesHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	db := data.New()
	token, err := db.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(token.Value)
	ctx := context.TODO()

	// Get list of machine sizes from DO client.
	sizes, err := client.ListSizes(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"sizes": sizes,
	})
}

func getCurrentMachinePricing(c *gin.Context) {
	db := data.New()
	token, err := db.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(token.Value)
	ctx := context.TODO()

	// Get the current machine pricing from DO client.
	pricing, err := client.GetCurrentPricing(ctx, WorkerTag)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"pricing": pricing,
	})
}

func listVPCsHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	db := data.New()
	token, err := db.Settings.GetSetting(types.DigitalOceanAccessToken)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	client, _ := machine.NewDigitalOceanClient(token.Value)
	ctx := context.TODO()

	// Get list of machine regions from DO client.
	vpcs, err := client.ListVPCs(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "machines not configured",
		})
		return
	}

	c.JSON(200, gin.H{
		"vpcs": vpcs,
	})
}
