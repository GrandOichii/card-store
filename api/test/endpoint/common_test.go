package endpoint_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/valkey-io/valkey-go"
	pgdb "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/router"
)

var (
	dbContainer    *postgres.PostgresContainer
	cacheContainer testcontainers.Container
)

func init() {
	var err error
	dbContainer, err = postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:latest"),
		postgres.WithDatabase("store-test-db"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "store-test-db",
			},
			Reuse: true,
		}),
	)
	if err != nil {
		panic(err)
	}
	err = dbContainer.Snapshot(context.Background())
	if err != nil {
		panic(err)
	}

	containerRequest := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:7.2.5",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
		Name:         "store-test-cache",
	}
	cacheContainer, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		panic(err)
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// TODO this is EXTREMLY slow
// each test has to create a separate container
// create a reusable db and cache that will be cleared after each test
func setupRouter(cardPageSize uint) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	testcontainers.Logger = log.New(&ioutils.NopWriter{}, "", 0)

	// db container
	ctx := context.Background()

	dbContainer.Restore(ctx)
	dbConn, err := dbContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		panic(err)
	}

	// cache container
	cacheConn, err := cacheContainer.Endpoint(context.Background(), "redis")
	if err != nil {
		panic(err)
	}

	// config
	config := config.Configuration{
		AuthKey: "test secret key",
		Port:    "8080",
		Db: config.DbConfiguration{
			ConnectionUri: dbConn,
			DbName:        "test_store",
			Cards: config.CardsDbConfiguration{
				PageSize: cardPageSize,
			},
		},
		Cache: config.CacheConfiguration{
			ConnectionUri: cacheConn,
		},
	}

	router := router.CreateRouter(&config)

	db, err := gorm.Open(pgdb.Open(config.Db.ConnectionUri), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cache, err := valkey.NewClient(valkey.MustParseURL(config.Cache.ConnectionUri))
	if err != nil {
		panic(err)
	}

	err = cache.Do(context.Background(), cache.
		B().
		Flushall().
		Build()).
		Error()
	if err != nil {
		panic(err)
	}

	return router, db
}

func toData(o interface{}) io.Reader {
	j, _ := json.Marshal(o)
	return bytes.NewBuffer(j)
}

func req(r *gin.Engine, t *testing.T, request string, path string, data interface{}, token string) (*httptest.ResponseRecorder, []byte) {
	var reqData io.Reader = nil
	if data != nil {
		reqData = toData(data)
	}
	req, err := http.NewRequest(request, path, reqData)
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	checkErr(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	result, err := io.ReadAll(w.Body)
	checkErr(t, err)
	return w, result
}

func createUser(r *gin.Engine, t *testing.T, username string, password string, email string) {
	req(r, t, "POST", "/api/v1/auth/register", dto.RegisterDetails{
		Password: password,
		Username: username,
		Email:    email,
	}, "")
}

func loginAs(r *gin.Engine, t *testing.T, username string, password string, email string) string {
	createUser(r, t, username, password, email)

	_, data := req(r, t, "POST", "/api/v1/auth/login", dto.LoginDetails{
		Username: username,
		Password: password,
	}, "")

	var res struct {
		Token string `json:"token"`
	}
	json.Unmarshal(data, &res)

	return res.Token
}

func createAdmin(r *gin.Engine, t *testing.T, db *gorm.DB) uint {
	username := "admin"
	createUser(r, t, username, "password", "admin@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
		Update("is_admin", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	var result model.User
	err = db.
		Where("username=?", username).
		Find(&result).
		Error

	if err != nil {
		t.Fatal(err)
	}

	return result.ID
}

func createCard(t *testing.T, db *gorm.DB, card *model.Card) uint {
	err := db.
		Create(card).
		Error
	if err != nil {
		t.Fatal(err)
	}
	var result model.Card
	err = db.
		Model(&model.Card{}).
		Where("name=?", card.Name).
		First(&result).
		Error
	if err != nil {
		t.Fatal(err)
	}

	return result.ID
}
