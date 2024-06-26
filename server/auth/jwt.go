package auth

import (
	"log"
	"net/http"
	"strconv"
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
	IDKey string = "id"
)

type JwtMiddleware struct {
	Middle                *jwt.GinJWTMiddleware
	AuthorizationCheckers []AuthorizationChecker
}

func NewJwtMiddleware(c *config.Configuration, authService service.AuthService, userRepo repository.UserRepository) *JwtMiddleware {
	result := new(JwtMiddleware)
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          c.JwtRealm,
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
					IDKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id := claims[IDKey].(string)
			userId, _ := strconv.ParseUint(id, 10, 32)
			user := userRepo.FindById(uint(userId))
			return user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals dto.LoginDetails
			if err := c.BindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			result, err := authService.Login(&loginVals)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return result, nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			user, ok := data.(*model.User)
			if !ok {
				return false
			}

			for _, checker := range result.AuthorizationCheckers {
				authorized, matches := checker.Check(c, user)
				if !matches {
					continue
				}
				return authorized
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.AbortWithStatusJSON(code, message)
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, cookie: token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie: token",

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

	result.Middle = authMiddleware
	return result
}
