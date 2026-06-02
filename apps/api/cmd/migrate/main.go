package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rentoutdoor/api/internal/infrastructure/config"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/migration"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/migration/seeds"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/mysql"
)

func main() {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := mysql.NewConnection(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer mysql.Close(db)

	migrator := migration.NewMigrator(db)

	switch command {
	case "up":
		log.Println("[MIGRATOR] Starting migrations up...")
		if err := migrator.MigrateUp(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("[MIGRATOR] Migrations up completed successfully.")

	case "down":
		log.Println("[MIGRATOR] Starting migration rollback...")
		if err := migrator.MigrateDown(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("[MIGRATOR] Migration rollback completed successfully.")

	case "seed":
		log.Println("[MIGRATOR] Starting database seeding...")
		if err := seeds.RunAll(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("[MIGRATOR] Database seeding completed successfully.")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down|seed]")
		os.Exit(1)
	}
}
