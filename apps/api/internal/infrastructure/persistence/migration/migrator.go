package migration

import (
	"embed"
	"fmt"
	"log"

	"gorm.io/gorm"
)

//go:embed *.sql
var migrationFS embed.FS

// Migrator handles database migrations using embedded SQL files.
type Migrator struct {
	db *gorm.DB
}

// MigrationRecord tracks applied migrations.
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Version   string `gorm:"type:varchar(100);not null;uniqueIndex"`
	AppliedAt string `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// NewMigrator creates a new migrator instance.
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Init creates the schema_migrations tracking table.
func (m *Migrator) Init() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// MigrateUp runs all pending up migrations in order.
func (m *Migrator) MigrateUp() error {
	if err := m.Init(); err != nil {
		return fmt.Errorf("failed to init migrations table: %w", err)
	}

	files, err := getMigrationFiles("up")
	if err != nil {
		return err
	}

	for _, file := range files {
		version := extractVersion(file)

		var count int64
		m.db.Model(&MigrationRecord{}).Where("version = ?", version).Count(&count)
		if count > 0 {
			continue
		}

		sql, err := migrationFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		log.Printf("[MIGRATION] Applying: %s", file)

		if err := m.db.Exec(string(sql)).Error; err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		m.db.Create(&MigrationRecord{Version: version})
		log.Printf("[MIGRATION] Applied: %s", file)
	}

	return nil
}

// MigrateDown rolls back the last applied migration.
func (m *Migrator) MigrateDown() error {
	if err := m.Init(); err != nil {
		return fmt.Errorf("failed to init migrations table: %w", err)
	}

	var last MigrationRecord
	if err := m.db.Order("id DESC").First(&last).Error; err != nil {
		return fmt.Errorf("no migrations to rollback: %w", err)
	}

	downFile := last.Version + ".down.sql"
	sql, err := migrationFS.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read down migration %s: %w", downFile, err)
	}

	log.Printf("[MIGRATION] Rolling back: %s", downFile)

	if err := m.db.Exec(string(sql)).Error; err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", downFile, err)
	}

	m.db.Where("version = ?", last.Version).Delete(&MigrationRecord{})
	log.Printf("[MIGRATION] Rolled back: %s", downFile)

	return nil
}

func getMigrationFiles(direction string) ([]string, error) {
	entries, err := migrationFS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var files []string
	suffix := "." + direction + ".sql"
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > len(suffix) {
			name := entry.Name()
			if name[len(name)-len(suffix):] == suffix {
				files = append(files, name)
			}
		}
	}
	return files, nil
}

func extractVersion(filename string) string {
	// "000001_create_users_table.up.sql" → "000001_create_users_table"
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			remainder := filename[:i]
			for j := len(remainder) - 1; j >= 0; j-- {
				if remainder[j] == '.' {
					return remainder[:j]
				}
			}
			return remainder
		}
	}
	return filename
}
