package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kasbon struct {
	ID              string     `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID          string     `gorm:"type:varchar(50);not null;index" json:"user_id"`
	Nominal         float64    `gorm:"type:decimal(15,2);not null" json:"nominal"`
	Reason          string     `gorm:"type:text;not null" json:"reason"`
	Status          string     `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	ApprovedBy      *string    `gorm:"type:varchar(50)" json:"approved_by"`
	ApprovedAt      *time.Time `json:"approved_at"`
	RejectionReason string     `gorm:"type:text" json:"rejection_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (k *Kasbon) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = uuid.New().String()
	}
	return nil
}