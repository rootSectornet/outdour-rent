package entity

type MountainDifficulty string

const (
	MountainDifficultyEasy     MountainDifficulty = "easy"
	MountainDifficultyModerate MountainDifficulty = "moderate"
	MountainDifficultyHard     MountainDifficulty = "hard"
	MountainDifficultyExpert   MountainDifficulty = "expert"
)

type Mountain struct {
	BaseModelNoSoftDelete
	Name            string             `gorm:"type:varchar(100);not null" json:"name"`
	Slug            string             `gorm:"type:varchar(120);not null;uniqueIndex:idx_mountains_slug" json:"slug"`
	Description     *string            `gorm:"type:text" json:"description,omitempty"`
	ElevationMeters *uint              `json:"elevation_meters,omitempty"`
	Difficulty      MountainDifficulty `gorm:"type:enum('easy','moderate','hard','expert');not null" json:"difficulty"`
	Location        *string            `gorm:"type:varchar(200)" json:"location,omitempty"`
	Latitude        *float64           `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude       *float64           `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	ImageURL        *string            `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	IsActive        bool               `gorm:"not null;default:true" json:"is_active"`

	// Relations
	EquipmentRecs []MountainEquipmentRec `gorm:"foreignKey:MountainID" json:"equipment_recs,omitempty"`
}

func (Mountain) TableName() string {
	return "mountains"
}

type RecPriority string

const (
	RecPriorityEssential   RecPriority = "essential"
	RecPriorityRecommended RecPriority = "recommended"
	RecPriorityOptional    RecPriority = "optional"
)

type MountainEquipmentRec struct {
	BaseModelCreateOnly
	MountainID string      `gorm:"type:char(36);not null;index:idx_recs_mountain_priority,priority:1;uniqueIndex:uidx_recs_mountain_category,priority:1" json:"mountain_id"`
	CategoryID string      `gorm:"type:char(36);not null;uniqueIndex:uidx_recs_mountain_category,priority:2" json:"category_id"`
	Priority   RecPriority `gorm:"type:enum('essential','recommended','optional');not null;index:idx_recs_mountain_priority,priority:2" json:"priority"`
	Notes      *string     `gorm:"type:varchar(500)" json:"notes,omitempty"`

	// Relations
	Mountain Mountain          `gorm:"foreignKey:MountainID;constraint:OnDelete:CASCADE" json:"mountain,omitempty"`
	Category EquipmentCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
}

func (MountainEquipmentRec) TableName() string {
	return "mountain_equipment_recs"
}
