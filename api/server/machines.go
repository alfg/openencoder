package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/machine"
	"github.com/gin-gonic/gin"
)

const settingKey = "DIGITAL_OCEAN_ACCESS_TOKEN"
const tag = "openencoder-worker"

type machineRequest struct {
	Provider string `json:"provider" binding:"required"`
	Region   string `json:"region" binding:"required"`
	Size     string `json:"size" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1,max=10"` // Max of 10 machines.
}

func machinesHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Get list of machines from DO client.
	machines, err := client.ListDropletByTag(ctx, tag)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"machines": machines,
	})
}

func createMachineHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

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

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Create machine.
	machine, err := client.CreateDroplets(ctx, json.Region, json.Size, json.Count)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Add resource to project?

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Create machine.
	machine, err := client.DeleteDropletByID(ctx, id)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"machine": machine,
	})
	return
}

func deleteMachineByTagHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Create machine.
	err := client.DeleteDropletByTag(ctx, tag)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"deleted": true,
	})
	return
}

func listMachineRegionsHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Get list of machine regions from DO client.
	regions, err := client.ListRegions(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"regions": regions,
	})
}

func listMachineSizesHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)

	// Role check.
	if !isAdminOrOperator(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Get list of machine sizes from DO client.
	sizes, err := client.ListSizes(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"sizes": sizes,
	})
}

func getCurrentMachinePricing(c *gin.Context) {
	d := data.New()
	token := d.Settings.GetSetting(settingKey).Value
	client, _ := machine.NewDigitalOceanClient(token)
	ctx := context.TODO()

	// Get the current machine pricing from DO client.
	pricing, err := client.GetCurrentPricing(ctx, tag)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"pricing": pricing,
	})
}
