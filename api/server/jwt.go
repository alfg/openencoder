package server

import (
	"errors"
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
		Realm:       JwtRealm,
		Key:         jwtKey,
		Timeout:     JwtTimeout,
		MaxRefresh:  JwtMaxRefresh,
		IdentityKey: JwtIdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*types.User); ok {
				return jwt.MapClaims{
					JwtIdentityKey: v.Username,
					JwtRoleKey:     v.Role,
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
				return nil, jwt.ErrFailedAuthentication
			}

			// Check the encrypted password.
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// Error with 403 if password needs to be reset.
			if user.ForcePasswordReset {
				return nil, errors.New("require password reset")
			}

			// Log-in the user.
			return &types.User{
				Username: user.Username,
				Role:     user.Role,
			}, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			// Only authorize if user has the following roles.
			if v, ok := data.(*types.User); ok &&
				(v.Role == "guest" || v.Role == "operator" || v.Role == "admin") {
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

		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(code, gin.H{
				"code":   code,
				"token":  message,
				"expire": time,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Error("JWT Error:" + err.Error())
	}
	return authMiddleware
}

func isAdminOrOperator(user interface{}) bool {
	role := user.(*types.User).Role
	if role != RoleOperator && role != RoleAdmin {
		return false
	}
	return true
}

func isOperator(user interface{}) bool {
	role := user.(*types.User).Role
	if role != RoleOperator {
		return false
	}
	return true
}

func isAdmin(user interface{}) bool {
	role := user.(*types.User).Role
	if role != RoleAdmin {
		return false
	}
	return true
}
