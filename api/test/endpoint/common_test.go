package endpoint_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pgdb "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"store.api/config"
	"store.api/dto"
	"store.api/router"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func setupRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	container, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:latest"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)

	if err != nil {
		panic(err)
	}
	conn, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		panic(err)
	}

	config := config.Configuration{
		AuthKey: "test secret key",
		Port:    "8080",
		Db: config.DbConfiguration{
			ConnectionUri: conn,
			DbName:        "test_store",
		},
	}

	router := router.CreateRouter(&config)

	db, err := gorm.Open(pgdb.Open(config.Db.ConnectionUri), &gorm.Config{})
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
