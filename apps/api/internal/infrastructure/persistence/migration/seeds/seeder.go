package seeds

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seeder defines the interface for all seeders.
type Seeder interface {
	Seed(db *gorm.DB) error
	Name() string
}

// RunAll executes all registered seeders in order.
func RunAll(db *gorm.DB) error {
	seeders := []Seeder{
		&AdminSeeder{},
		&CategorySeeder{},
		&MountainSeeder{},
	}

	for _, s := range seeders {
		log.Printf("[SEED] Running: %s", s.Name())
		if err := s.Seed(db); err != nil {
			return err
		}
		log.Printf("[SEED] Completed: %s", s.Name())
	}

	return nil
}

// newUUID generates a new UUID v4 string.
func newUUID() string {
	return uuid.New().String()
}
