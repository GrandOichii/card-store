package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sethvargo/go-envconfig"
)

type StoreConfiguration struct {
}

type CardsDbConfiguration struct {
	PageSize uint `json:"pageSize" env:"PAGE_SIZE"`
}

type DbConfiguration struct {
	ConnectionUri string               `json:"connectionUri" env:"CONNECTION_URI"`
	DbName        string               `json:"dbName" env:"NAME"`
	Cards         CardsDbConfiguration `json:"cards" env:",prefix=CARDS_"`
}

type CacheConfiguration struct {
	ConnectionUri string `json:"connectionUri" env:"CONNECTION_URI"`
}

type Configuration struct {
	Host     string             `json:"host" env:"HOST"`
	Port     string             `json:"port" env:"PORT,required"`
	Db       DbConfiguration    `json:"db" env:",prefix=DB_"`
	Cache    CacheConfiguration `json:"cache" env:",prefix=CACHE_"`
	Store    StoreConfiguration `json:"store" env:",prefix=STORE_"`
	AuthKey  string             `json:"authKey" env:"AUTH_KEY"`
	JwtRealm string             `json:"jwtRealm" env:"JWT_REALM"`
}

func ReadConfig(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	result := &Configuration{}
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ReadEnvConfig() (*Configuration, error) {
	ctx := context.Background()

	var result Configuration
	if err := envconfig.Process(ctx, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
