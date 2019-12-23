package server

import (
	"fmt"
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

type userUpdateRequest struct {
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"`
	VerifyPassword  string `json:"verify_password"`
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

func getUserHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	username := user.(*types.User).Username

	db := data.New()
	u, _ := db.Users.GetUserByUsername(username)

	c.JSON(200, gin.H{
		"status":   http.StatusOK,
		"username": u.Username,
		"role":     u.Role,
	})
}

func updateUserHandler(c *gin.Context) {
	user, _ := c.Get(identityKey)
	username := user.(*types.User).Username
	// if username != "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"message": "unauthorized",
	// 	})
	// }

	// Decode json.
	var json userUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from DB.
	db := data.New()
	u, _ := db.Users.GetUserByUsername(username)

	// Verify user credentials.
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(json.CurrentPassword))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "unauthorized",
		})
		return
	}

	// Update user struct with data provided.
	if json.Username != "" {
		u.Username = json.Username
	}

	// Create new password hash if new_password is provided.
	if json.NewPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(json.NewPassword), bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "error updating user",
			})
			return
		}
		u.Password = string(hash)
	}

	u, err = db.Users.UpdateUserByID(u.ID, u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "error updating user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    username,
		"message": "user updated",
	})
}
