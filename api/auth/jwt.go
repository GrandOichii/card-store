package auth

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/service"
)

const (
	IDKey string = "username"
)

type JwtMiddleware struct {
	Middle *jwt.GinJWTMiddleware
}

func NewJwtMiddleware(config *config.Configuration, userService service.UserService) *JwtMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(config.AuthKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,

		IdentityKey: IDKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*dto.PrivateUserInfo); ok {
				return jwt.MapClaims{
					IDKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			fmt.Printf("claims: %v\n", claims)
			return &model.User{
				Username: claims[IDKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals dto.LoginDetails
			if err := c.BindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			result, err := userService.Login(&loginVals)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return result, nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*models.User); ok && v.UserName == "admin" {
			// 	return true
			// }

			// TODO

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error": message,
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
		log.Fatal(err.Error())
	}

	err = authMiddleware.MiddlewareInit()

	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	result := &JwtMiddleware{
		Middle: authMiddleware,
	}
	return result
}
