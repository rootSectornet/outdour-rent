package entity

// AllModels returns all GORM models for auto-migration (development only).
// Production uses SQL migration files.
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&UserSession{},
		&PasswordReset{},
		&Store{},
		&StoreDocument{},
		&EquipmentCategory{},
		&Equipment{},
		&EquipmentPhoto{},
		&EquipmentPricing{},
		&EquipmentMaintenance{},
		&InventoryReservation{},
		&ReservationDateLock{},
		&Order{},
		&OrderItem{},
		&OrderStatusHistory{},
		&Payment{},
		&Deposit{},
		&Refund{},
		&Review{},
		&Notification{},
		&Mountain{},
		&MountainEquipmentRec{},
		&AdminAuditLog{},
	}
}
