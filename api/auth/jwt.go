package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/service"
)

const (
	IDKey string = "username"
)

type JwtMiddleware struct {
	Middle *jwt.GinJWTMiddleware
}

func NewJwtMiddleware(c *config.Configuration, userService service.UserService, userRepo repository.UserRepository) *JwtMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "test zone", // TODO
		Key:            []byte(c.AuthKey),
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		SendCookie:     true,
		SecureCookie:   false, // ! non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		CookieDomain:   c.Host + ":" + c.Port,
		CookieName:     "token",                  // default jwt
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode

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
			username := claims[IDKey].(string)
			user := userRepo.FindByUsername(username)
			return user
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
			user, ok := data.(*model.User)
			if !ok {
				return false
			}
			// c.Request.

			// TODO this can be done better
			if c.Request.Method != "GET" && c.Request.URL.Path == "/api/v1/card" {
				fmt.Printf("user: %v\n", user)
				fmt.Printf("user.Email: %v\n", user.Email)
				fmt.Printf("user.IsAdmin: %v\n", user.IsAdmin)
				fmt.Printf("user.Verified: %v\n", user.Verified)
				return user.IsAdmin && user.Verified
			}

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
