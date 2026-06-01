package seeds

import (
	"time"

	"gorm.io/gorm"
)

// MountainSeeder seeds popular Indonesian mountains.
type MountainSeeder struct{}

func (s *MountainSeeder) Name() string {
	return "MountainSeeder"
}

func (s *MountainSeeder) Seed(db *gorm.DB) error {
	type Mountain struct {
		ID              string   `gorm:"type:char(36);primaryKey"`
		Name            string   `gorm:"type:varchar(100)"`
		Slug            string   `gorm:"type:varchar(120)"`
		Description     *string  `gorm:"type:text"`
		ElevationMeters *uint    `gorm:"column:elevation_meters"`
		Difficulty      string   `gorm:"type:enum('easy','moderate','hard','expert')"`
		Location        *string  `gorm:"type:varchar(200)"`
		Latitude        *float64 `gorm:"type:decimal(10,8)"`
		Longitude       *float64 `gorm:"type:decimal(11,8)"`
		IsActive        bool
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}

	var count int64
	db.Table("mountains").Count(&count)
	if count > 0 {
		return nil
	}

	now := time.Now()

	ptr := func(s string) *string { return &s }
	ptrUint := func(u uint) *uint { return &u }
	ptrFloat := func(f float64) *float64 { return &f }

	mountains := []Mountain{
		{
			ID: newUUID(), Name: "Gunung Semeru", Slug: "gunung-semeru",
			Description: ptr("The highest mountain in Java at 3,676m. Known as Mahameru."),
			ElevationMeters: ptrUint(3676), Difficulty: "hard",
			Location: ptr("East Java"), Latitude: ptrFloat(-8.1080), Longitude: ptrFloat(112.9220),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Rinjani", Slug: "gunung-rinjani",
			Description: ptr("Second highest volcano in Indonesia at 3,726m with a stunning crater lake."),
			ElevationMeters: ptrUint(3726), Difficulty: "hard",
			Location: ptr("West Nusa Tenggara"), Latitude: ptrFloat(-8.4110), Longitude: ptrFloat(116.4570),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Prau", Slug: "gunung-prau",
			Description: ptr("Popular beginner mountain in Central Java with golden sunrise views."),
			ElevationMeters: ptrUint(2565), Difficulty: "easy",
			Location: ptr("Central Java"), Latitude: ptrFloat(-7.1850), Longitude: ptrFloat(109.9220),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Merbabu", Slug: "gunung-merbabu",
			Description: ptr("Stratovolcano in Central Java with beautiful savanna landscapes."),
			ElevationMeters: ptrUint(3145), Difficulty: "moderate",
			Location: ptr("Central Java"), Latitude: ptrFloat(-7.4550), Longitude: ptrFloat(110.4400),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Kerinci", Slug: "gunung-kerinci",
			Description: ptr("Highest volcano in Sumatra at 3,805m. Located in Kerinci Seblat National Park."),
			ElevationMeters: ptrUint(3805), Difficulty: "hard",
			Location: ptr("Jambi, Sumatra"), Latitude: ptrFloat(-1.6970), Longitude: ptrFloat(101.2640),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Gede Pangrango", Slug: "gunung-gede-pangrango",
			Description: ptr("Twin volcano near Jakarta. Popular weekend hiking destination."),
			ElevationMeters: ptrUint(2958), Difficulty: "moderate",
			Location: ptr("West Java"), Latitude: ptrFloat(-6.7850), Longitude: ptrFloat(106.9830),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Bromo", Slug: "gunung-bromo",
			Description: ptr("Iconic active volcano in East Java. Part of Tengger caldera."),
			ElevationMeters: ptrUint(2329), Difficulty: "easy",
			Location: ptr("East Java"), Latitude: ptrFloat(-7.9425), Longitude: ptrFloat(112.9530),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: newUUID(), Name: "Gunung Lawu", Slug: "gunung-lawu",
			Description: ptr("Stratovolcano on the border of Central and East Java."),
			ElevationMeters: ptrUint(3265), Difficulty: "moderate",
			Location: ptr("Central/East Java"), Latitude: ptrFloat(-7.6250), Longitude: ptrFloat(111.1920),
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
	}

	return db.Table("mountains").Create(&mountains).Error
}
