package server

import (
	"fmt"
	"log"
	"time"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/data"
	"github.com/alfg/openencoder/api/helpers"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

const (
	identityKey = "id"
	roleKey     = "role"
)

var jwtKey []byte

// User demo
type User struct {
	Username string
	Role     string
}

func jwtMiddleware() *jwt.GinJWTMiddleware {

	// Set the JWT Key if provided in config. Otherwise, generate a random one.
	key := config.Get().JWTKey
	if key == "" {
		jwtKey = helpers.GenerateRandomKey(16)
	} else {
		jwtKey = []byte(key)
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "openencoder",
		Key:         jwtKey,
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					roleKey:     v.Role,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
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

			user, err := data.GetUserByUsername(userID)
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
			return &User{
				Username: user.Username,
				Role:     user.Role,
			}, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			// Only authorize if user is an operator.
			if v, ok := data.(*User); ok && v.Role == "operator" {
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
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
