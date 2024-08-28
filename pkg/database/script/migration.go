package script

import (
	"context"
	"go-meechok/config"
	migration "go-meechok/pkg/database/migrations"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// Initialize config
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	migration.AuthMigrate(ctx, cfg)
	migration.UserMigrate(ctx, cfg)

}
