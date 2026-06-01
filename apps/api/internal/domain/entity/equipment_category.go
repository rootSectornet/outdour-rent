package entity

type EquipmentCategory struct {
	BaseModel
	ParentID  *string `gorm:"type:char(36);index:idx_categories_parent_id" json:"parent_id,omitempty"`
	Name      string  `gorm:"type:varchar(100);not null" json:"name"`
	Slug      string  `gorm:"type:varchar(120);not null;uniqueIndex:idx_categories_slug" json:"slug"`
	IconURL   *string `gorm:"type:varchar(500)" json:"icon_url,omitempty"`
	SortOrder int     `gorm:"not null;default:0" json:"sort_order"`
	IsActive  bool    `gorm:"not null;default:true" json:"is_active"`

	// Relations
	Parent   *EquipmentCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []EquipmentCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (EquipmentCategory) TableName() string {
	return "equipment_categories"
}
