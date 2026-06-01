package entity

import (
	"time"
)

type StoreStatus string

const (
	StoreStatusPending   StoreStatus = "pending"
	StoreStatusVerified  StoreStatus = "verified"
	StoreStatusSuspended StoreStatus = "suspended"
	StoreStatusRejected  StoreStatus = "rejected"
)

type Store struct {
	BaseModel
	OwnerID     string      `gorm:"type:char(36);not null;index:idx_stores_owner_id" json:"owner_id"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string      `gorm:"type:varchar(120);not null;uniqueIndex:idx_stores_slug" json:"slug"`
	Description *string     `gorm:"type:text" json:"description,omitempty"`
	Phone       string      `gorm:"type:varchar(20);not null" json:"phone"`
	Email       string      `gorm:"type:varchar(255);not null" json:"email"`
	Address     string      `gorm:"type:text;not null" json:"address"`
	City        string      `gorm:"type:varchar(100);not null;index:idx_stores_city" json:"city"`
	Province    string      `gorm:"type:varchar(100);not null" json:"province"`
	PostalCode  *string     `gorm:"type:varchar(10)" json:"postal_code,omitempty"`
	Latitude    *float64    `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude   *float64    `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	LogoURL     *string     `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
	BannerURL   *string     `gorm:"type:varchar(500)" json:"banner_url,omitempty"`
	Status      StoreStatus `gorm:"type:enum('pending','verified','suspended','rejected');not null;default:'pending';index:idx_stores_status" json:"status"`
	VerifiedAt  *time.Time  `json:"verified_at,omitempty"`
	RatingAvg   float64     `gorm:"type:decimal(3,2);not null;default:0.00" json:"rating_avg"`
	RatingCount uint        `gorm:"not null;default:0" json:"rating_count"`

	// Relations
	Owner      User            `gorm:"foreignKey:OwnerID;constraint:OnDelete:RESTRICT" json:"owner,omitempty"`
	Equipment  []Equipment     `gorm:"foreignKey:StoreID" json:"equipment,omitempty"`
	Orders     []Order         `gorm:"foreignKey:StoreID" json:"orders,omitempty"`
	Documents  []StoreDocument `gorm:"foreignKey:StoreID" json:"documents,omitempty"`
}

func (Store) TableName() string {
	return "stores"
}

type DocumentType string

const (
	DocumentTypeKTP   DocumentType = "ktp"
	DocumentTypeNPWP  DocumentType = "npwp"
	DocumentTypeSIU   DocumentType = "siu"
	DocumentTypeOther DocumentType = "other"
)

type StoreDocument struct {
	BaseModelCreateOnly
	StoreID        string       `gorm:"type:char(36);not null;index:idx_store_docs_store_id" json:"store_id"`
	DocumentType   DocumentType `gorm:"type:enum('ktp','npwp','siu','other');not null" json:"document_type"`
	DocumentURL    string       `gorm:"type:varchar(500);not null" json:"document_url"`
	VerifiedAt     *time.Time   `json:"verified_at,omitempty"`
	RejectedReason *string      `gorm:"type:varchar(500)" json:"rejected_reason,omitempty"`

	// Relations
	Store Store `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE" json:"store,omitempty"`
}

func (StoreDocument) TableName() string {
	return "store_documents"
}
