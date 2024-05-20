package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"store.api/auth"
	"store.api/config"
	"store.api/controller"
	"store.api/model"
	"store.api/repository"
	"store.api/service"
	"store.api/utility"

	gin "github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
)

func CreateRouter(config *config.Configuration) *gin.Engine {
	result := gin.Default()

	result.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// TODO
			return true
			// return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// database
	dbClient, err := dbConnect(config)
	if err != nil {
		panic(err)
	}

	err = dbConfig(dbClient)

	if err != nil {
		panic(err)
	}

	// repositories
	userRepo := repository.NewUserDbRepository(dbClient, config)
	cardRepo := repository.NewCardDbRepository(dbClient, config)

	configRouter(result, config, userRepo, cardRepo)

	result.GET("/api/v1/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hi!")
	})

	return result
}

func configRouter(router *gin.Engine, config *config.Configuration, userRepo repository.UserRepository, cardRepo repository.CardRepository) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// services
	userService := service.NewUserServiceImpl(
		userRepo,
		validate,
	)
	cardService := service.NewCardServiceImpl(
		cardRepo,
		userRepo,
		validate,
	)

	// middleware
	auth := auth.NewJwtMiddleware(config, userService, userRepo)

	// controllers
	api := router.Group("/api/v1")

	cardController := controller.NewCardController(
		cardService,
		auth.Middle.MiddlewareFunc(),
		utility.Extract,
	)
	cardController.Configure(api)

	authController := controller.NewAuthController(
		userService,
		auth.Middle.LoginHandler,
	)
	authController.Configure(api)
}

func dbConnect(config *config.Configuration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Db.ConnectionUri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dbConfig(db *gorm.DB) error {

	err := db.AutoMigrate(
		&model.User{},
		&model.Card{},
	)
	if err != nil {
		return err
	}

	return nil
}
