package server

import (
	"github.com/alfg/openencoder/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"name":    "openencoder",
		"version": "0.0.1",
		"github":  "https://github.com/alfg/openencoder",
		"user_id": claims["id"],
		"user":    user.(*types.User).Username,
		"role":    user.(*types.User).Role,
	})
}
