package entity

type AdminAuditLog struct {
	BaseModelCreateOnly
	AdminID    string  `gorm:"type:char(36);not null;index:idx_audit_admin_created,priority:1" json:"admin_id"`
	Action     string  `gorm:"type:varchar(50);not null;index:idx_audit_action_created,priority:1" json:"action"`
	EntityType string  `gorm:"type:varchar(50);not null;index:idx_audit_entity,priority:1" json:"entity_type"`
	EntityID   string  `gorm:"type:char(36);not null;index:idx_audit_entity,priority:2" json:"entity_id"`
	OldValues  *string `gorm:"type:json" json:"old_values,omitempty"`
	NewValues  *string `gorm:"type:json" json:"new_values,omitempty"`
	IPAddress  *string `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string `gorm:"type:varchar(500)" json:"user_agent,omitempty"`

	// Relations
	Admin User `gorm:"foreignKey:AdminID;constraint:OnDelete:RESTRICT" json:"admin,omitempty"`
}

func (AdminAuditLog) TableName() string {
	return "admin_audit_logs"
}
