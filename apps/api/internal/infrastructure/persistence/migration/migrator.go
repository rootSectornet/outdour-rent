package migration

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

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
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Version   string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	AppliedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
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

		if err := m.executeSQLScript(string(sql)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		if err := m.db.Create(&MigrationRecord{Version: version}).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}
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

	if err := m.executeSQLScript(string(sql)); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", downFile, err)
	}

	if err := m.db.Where("version = ?", last.Version).Delete(&MigrationRecord{}).Error; err != nil {
		return fmt.Errorf("failed to delete migration record %s: %w", downFile, err)
	}
	log.Printf("[MIGRATION] Rolled back: %s", downFile)

	return nil
}

func (m *Migrator) executeSQLScript(script string) error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to access sql.DB: %w", err)
	}

	ctx := context.Background()
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire database connection: %w", err)
	}
	defer conn.Close()

	for _, statement := range splitSQLStatements(script) {
		if _, err := conn.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("statement %q failed: %w", previewStatement(statement), err)
		}
	}

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

func splitSQLStatements(script string) []string {
	var statements []string
	var current strings.Builder

	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(script); i++ {
		ch := script[i]
		next := byte(0)
		if i+1 < len(script) {
			next = script[i+1]
		}

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			}
			continue
		}

		if inBlockComment {
			if ch == '*' && next == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if !inSingleQuote && !inDoubleQuote && !inBacktick {
			if isLineCommentStart(script, i) {
				inLineComment = true
				i++
				continue
			}
			if ch == '/' && next == '*' {
				inBlockComment = true
				i++
				continue
			}
			if ch == ';' {
				statement := strings.TrimSpace(current.String())
				if statement != "" {
					statements = append(statements, statement)
				}
				current.Reset()
				continue
			}
		}

		switch ch {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				if inSingleQuote && next == '\'' {
					current.WriteByte(ch)
					current.WriteByte(next)
					i++
					continue
				}
				if !isEscaped(script, i) {
					inSingleQuote = !inSingleQuote
				}
			}
		case '"':
			if !inSingleQuote && !inBacktick {
				if inDoubleQuote && next == '"' {
					current.WriteByte(ch)
					current.WriteByte(next)
					i++
					continue
				}
				if !isEscaped(script, i) {
					inDoubleQuote = !inDoubleQuote
				}
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		}

		current.WriteByte(ch)
	}

	if statement := strings.TrimSpace(current.String()); statement != "" {
		statements = append(statements, statement)
	}

	return statements
}

func isLineCommentStart(script string, idx int) bool {
	if idx+2 >= len(script) || script[idx] != '-' || script[idx+1] != '-' {
		return false
	}

	next := script[idx+2]
	if next != ' ' && next != '\t' && next != '\r' && next != '\n' {
		return false
	}

	if idx == 0 {
		return true
	}

	prev := script[idx-1]
	return prev == '\n' || prev == '\r' || prev == ' ' || prev == '\t'
}

func isEscaped(script string, idx int) bool {
	backslashes := 0
	for i := idx - 1; i >= 0 && script[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}

func previewStatement(statement string) string {
	const maxLen = 80

	normalized := strings.Join(strings.Fields(statement), " ")
	if len(normalized) <= maxLen {
		return normalized
	}

	return normalized[:maxLen] + "..."
}
