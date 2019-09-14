package server

import (
	"fmt"
	"log"
	"time"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// JWT settings.
const (
	realm       = "openencoder"
	identityKey = "id"
	roleKey     = "role"
	timeout     = time.Hour // Duration a JWT is valid.
	maxRefresh  = time.Hour // Duration a JWT can be refreshed.
)

var jwtKey []byte

func jwtMiddleware() *jwt.GinJWTMiddleware {

	// Set the JWT Key if provided in config. Otherwise, generate a random one.
	key := config.Get().JWTKey
	if key == "" {
		jwtKey = helpers.GenerateRandomKey(16)
	} else {
		jwtKey = []byte(key)
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       realm,
		Key:         jwtKey,
		Timeout:     timeout,
		MaxRefresh:  maxRefresh,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*types.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					roleKey:     v.Role,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &types.User{
				Username: claims["id"].(string),
				Role:     claims["role"].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			db := data.New()
			user, err := db.Users.GetUserByUsername(userID)
			if err != nil {
				fmt.Println(err)
				return nil, jwt.ErrFailedAuthentication
			}

			// Check the encrypted password.
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				fmt.Println(err)
				return nil, jwt.ErrFailedAuthentication
			}

			// Log-in the user.
			return &types.User{
				Username: user.Username,
				Role:     user.Role,
			}, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			// Only authorize if user is an operator.
			if v, ok := data.(*types.User); ok && v.Role == "operator" {
				return true
			}
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
