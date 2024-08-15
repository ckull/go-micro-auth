package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go-auth/utils"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type (
	Config struct {
		*Server
		*Db
		*Jwt
		*Grpc
		*Facebook
		*Google
	}

	Server struct {
		Port int64
	}

	Db struct {
		URI string
	}

	Facebook struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
		Endpoint     string
		Scopes       []string
		*oauth2.Config
	}

	Google struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
		Endpoint     string
		Scopes       []string
		*oauth2.Config
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
			AccessTokenDuration:  utils.ParseStringToInt(os.Getenv("ACCESS_TOKEN_DURATION")),
			RefreshTokenDuration: utils.ParseStringToInt(os.Getenv("REFRESH_TOKEN_DURATION")),
			ApiDuration:          utils.ParseStringToInt(os.Getenv("ACCESS_TOKEN_DURATION")),
		},
		Facebook: &Facebook{
			ClientID:     os.Getenv("OAUTH2_FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH2_FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH2_FACEBOOK_REDIRECT_URL"),
			Endpoint:     os.Getenv("OAUTH2_FACEBOOK_ENDPOINT"),
			Scopes: func() []string {
				var scopes []string
				jsonString := os.Getenv("OAUTH2_FACEBOOK_SCOPES")
				err := json.Unmarshal([]byte(jsonString), &scopes)
				if err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
					return nil
				}
				return scopes
			}(),
		},
		Google: &Google{
			ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL"),
			Endpoint:     os.Getenv("OAUTH2_GOOGLE_ENDPOINT"),
			Scopes: func() []string {
				var scopes []string
				jsonString := os.Getenv("OAUTH2_GOOGLE_SCOPES")
				err := json.Unmarshal([]byte(jsonString), &scopes)
				if err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
					return nil
				}
				return scopes
			}(),
		},
	}
}
