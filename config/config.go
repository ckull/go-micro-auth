package config

import (
	"log"
	"os"

	"go-auth/utils"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		*Server
		*Db
		*Jwt
		*Grpc
	}

	Server struct {
		Port int64
	}

	Db struct {
		URI string
	}

	Jwt struct {
		AccessTokenSecret    string
		RefreshTokenSecret   string
		ApiSecret            string
		AccessTokenDuration  int64
		RefreshTokenDuration int64
		ApiDuration          int64
	}

	Grpc struct {
		AuthUrl string
		UserUrl string
	}
)

func LoadConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Server: &Server{
			Port: utils.ParseStringToInt(os.Getenv("SERVER_PORT")),
		},
		Db: &Db{
			URI: os.Getenv("DB_URI"),
		},
		Jwt: &Jwt{
			AccessTokenSecret:    os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshTokenSecret:   os.Getenv("REFRESH_TOKEN_SECRET"),
			ApiSecret:            os.Getenv("API_SECRET"),
			AccessTokenDuration:  utils.ParseStringToInt("ACCESS_TOKEN_DURATION"),
			RefreshTokenDuration: utils.ParseStringToInt("ACCESS_TOKEN_DURATION"),
			ApiDuration:          utils.ParseStringToInt("ACCESS_TOKEN_SECRET"),
		},
	}
}
