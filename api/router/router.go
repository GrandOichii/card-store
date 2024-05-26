package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
	"store.api/auth"
	"store.api/cache"
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

	// cache
	cacheClient, err := cacheConnect(config)
	if err != nil {
		panic(err)
	}

	// repositories
	userRepo := repository.NewUserDbRepository(dbClient, config)
	cardRepo := repository.NewCardDbRepository(
		dbClient,
		config,
		cache.NewCardValkeyCache(cacheClient),
	)
	collectionRepo := repository.NewCollectionDbRepository(
		dbClient,
		config,
		cache.NewCollectionValkeyCache(cacheClient),
	)

	configRouter(result, config, cacheClient, userRepo, cardRepo, collectionRepo)

	return result
}

func configRouter(
	router *gin.Engine,
	config *config.Configuration,
	cacheClient valkey.Client,
	userRepo repository.UserRepository,
	cardRepo repository.CardRepository,
	collectionRepo repository.CollectionRepository,
) {
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
	collectionService := service.NewCollectionServiceImpl(
		collectionRepo,
		userRepo,
		validate,
	)

	// middleware
	authentication := auth.NewJwtMiddleware(
		config,
		userService,
		userRepo,
	)

	// controllers
	cardController := controller.NewCardController(
		cardService,
		authentication.Middle.MiddlewareFunc(),
		utility.Extract,
	)

	authController := controller.NewAuthController(
		userService,
		authentication.Middle.LoginHandler,
	)

	userController := controller.NewUserController(
		userService,
		authentication.Middle.MiddlewareFunc(),
	)

	collectionController := controller.NewCollectionController(
		collectionService,
		authentication.Middle.MiddlewareFunc(),
		utility.Extract,
	)

	api := router.Group("/api/v1")
	controllers := []controller.Controller{
		cardController,
		authController,
		userController,
		collectionController,
	}
	for _, c := range controllers {
		c.ConfigureApi(api)
	}

	authentication.AuthorizationCheckers = []auth.AuthorizationChecker{
		cardController,
		userController,
		collectionController,
	}
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
		&model.CardType{},
		&model.Language{},
		&model.Collection{},
		&model.CollectionSlot{},
	)
	if err != nil {
		return err
	}

	return nil
}

func cacheConnect(config *config.Configuration) (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.MustParseURL(config.Cache.ConnectionUri))
	if err != nil {
		return nil, err
	}
	return client, err
}
