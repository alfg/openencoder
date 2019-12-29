package server

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userProfileUpdateRequest struct {
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"`
	VerifyPassword  string `json:"verify_password"`
}

type userUpdateRequest struct {
	Active bool   `json:"active"`
	Role   string `json:"role" binding:"eq=admin|eq=operator|eq=guest|eq="`
}

type userPasswordUpdateRequest struct {
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	VerifyPassword  string `json:"verify_password"`
}

// registerHandler handles the request to register a new user.
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

// getUserProfileHandler handles the request to get the current user profile data.
func getUserProfileHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)
	username := user.(*types.User).Username

	db := data.New()
	u, _ := db.Users.GetUserByUsername(username)

	c.JSON(200, gin.H{
		"status":   http.StatusOK,
		"username": u.Username,
		"role":     u.Role,
	})
}

// updateUserProfileHandler handles the request to update the current user profile data.
func updateUserProfileHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)
	username := user.(*types.User).Username

	// Decode json.
	var json userProfileUpdateRequest
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

	// Update the user.
	u, err = db.Users.UpdateUserByID(int(u.ID), u)
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

// updatePasswordHandler handles the request to update a user password if the user
// credentials are validated.
func updatePasswordHandler(c *gin.Context) {
	// Decode json.
	var json userProfileUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "validation error. please provide username, current password and new password",
		})
		return
	}

	// Get user from DB.
	db := data.New()
	u, _ := db.Users.GetUserByUsername(json.Username)

	// Verify user credentials.
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(json.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "unauthorized",
		})
		return
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

	// Update the user.
	u, err = db.Users.UpdateUserPasswordByID(u.ID, u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "error updating user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    u.Username,
		"message": "user updated",
	})
}

// getUsersHandler handles the request to get all users for user management.
func getUsersHandler(c *gin.Context) {
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdmin(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	page := c.DefaultQuery("page", "1")
	count := c.DefaultQuery("count", "10")
	pageInt, _ := strconv.Atoi(page)
	countInt, _ := strconv.Atoi(count)

	if page == "0" {
		pageInt = 1
	}

	var wg sync.WaitGroup
	var users *[]types.User
	var usersCount int

	db := data.New()
	wg.Add(1)
	go func() {
		users = db.Users.GetUsers((pageInt-1)*countInt, countInt)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		usersCount = db.Users.GetUsersCount()
		wg.Done()
	}()
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"users": users,
		"count": usersCount,
	})
}

// updateUserByIDHandler handles the request update a user for user management.
func updateUserByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, _ := c.Get(JwtIdentityKey)

	// Role check.
	if !isAdmin(user) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Decode json.
	var json userUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := data.New()
	u, err := db.Users.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User does not exist",
		})
		return
	}

	// Disallow updates on master user.
	if id != 1 {
		// Set role.
		if json.Role != "" {
			u.Role = json.Role
		}

		// Set active status.
		u.Active = json.Active
	}

	updatedUser, _ := db.Users.UpdateUserByID(id, u)
	c.JSON(http.StatusOK, updatedUser)
}
