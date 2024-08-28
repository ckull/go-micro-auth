package config

import (
	"log"
	"os"

	"go-meechok/utils"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type (
	Config struct {
		*Server
		*Db
		*Jwt
		*Grpc
		Facebook *oauth2.Config
		Google   *oauth2.Config
		*Redis
		*Kafka
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

	Redis struct {
		Address  string
		Password string
		DB       int
	}

	Grpc struct {
		AuthUrl      string
		UserUrl      string
		InventoryUrl string
	}

	Kafka struct {
		Brokers []string
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
		Grpc: &Grpc{
			AuthUrl:      ":50051",
			UserUrl:      ":50052",
			InventoryUrl: ":50053",
		},
		Facebook: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH2_FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH2_FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH2_FACEBOOK_REDIRECT_URL"),
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"email"},
		},
		Google: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH2_GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"email"},
		},
		Kafka: &Kafka{
			Brokers: []string{""},
		},
	}
}
