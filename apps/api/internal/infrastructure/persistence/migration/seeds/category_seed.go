package seeds

import (
	"time"

	"gorm.io/gorm"
)

// CategorySeeder seeds default equipment categories.
type CategorySeeder struct{}

func (s *CategorySeeder) Name() string {
	return "CategorySeeder"
}

func (s *CategorySeeder) Seed(db *gorm.DB) error {
	type Category struct {
		ID        string  `gorm:"type:char(36);primaryKey"`
		ParentID  *string `gorm:"type:char(36)"`
		Name      string  `gorm:"type:varchar(100)"`
		Slug      string  `gorm:"type:varchar(120)"`
		SortOrder int
		IsActive  bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	var count int64
	db.Table("equipment_categories").Count(&count)
	if count > 0 {
		return nil
	}

	now := time.Now()

	// Root categories
	categories := []Category{
		{ID: newUUID(), Name: "Shelter", Slug: "shelter", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Sleeping", Slug: "sleeping", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Cooking", Slug: "cooking", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Clothing", Slug: "clothing", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Navigation", Slug: "navigation", SortOrder: 5, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Lighting", Slug: "lighting", SortOrder: 6, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Backpacks", Slug: "backpacks", SortOrder: 7, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Trekking Poles", Slug: "trekking-poles", SortOrder: 8, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Safety & First Aid", Slug: "safety-first-aid", SortOrder: 9, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), Name: "Electronics", Slug: "electronics", SortOrder: 10, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	// Sub-categories for Shelter
	shelterID := categories[0].ID
	subCategories := []Category{
		{ID: newUUID(), ParentID: &shelterID, Name: "Tents", Slug: "tents", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &shelterID, Name: "Tarps", Slug: "tarps", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &shelterID, Name: "Flysheets", Slug: "flysheets", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	// Sub-categories for Sleeping
	sleepingID := categories[1].ID
	sleepingSubs := []Category{
		{ID: newUUID(), ParentID: &sleepingID, Name: "Sleeping Bags", Slug: "sleeping-bags", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &sleepingID, Name: "Sleeping Pads", Slug: "sleeping-pads", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &sleepingID, Name: "Hammocks", Slug: "hammocks", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	// Sub-categories for Cooking
	cookingID := categories[2].ID
	cookingSubs := []Category{
		{ID: newUUID(), ParentID: &cookingID, Name: "Stoves", Slug: "stoves", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &cookingID, Name: "Cookware", Slug: "cookware", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: newUUID(), ParentID: &cookingID, Name: "Water Filters", Slug: "water-filters", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	allCategories := append(categories, subCategories...)
	allCategories = append(allCategories, sleepingSubs...)
	allCategories = append(allCategories, cookingSubs...)

	return db.Table("equipment_categories").Create(&allCategories).Error
}
