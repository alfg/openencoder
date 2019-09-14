package server

import (
	"net/http"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func registerHandler(c *gin.Context) {
	// Decode json.
	var json registerRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error creating user",
		})
		return
	}

	// Create Job and push the work to work queue.
	user := types.User{
		Username: json.Username,
		Password: string(hash),
		Role:     "guest",
	}

	db := data.New()
	u, err := db.Users.CreateUser(user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error creating user",
		})
		return
	}

	c.JSON(200, gin.H{
		"user":    u.Username,
		"message": "user created",
	})
}
