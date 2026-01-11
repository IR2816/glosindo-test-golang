package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Presensi struct {
	ID              string     `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID          string     `gorm:"type:varchar(50);not null" json:"user_id"`
	Date            time.Time  `gorm:"type:date;not null" json:"date"`
	CheckInTime     *time.Time `gorm:"type:timestamp" json:"check_in_time"`
	CheckOutTime    *time.Time `gorm:"type:timestamp" json:"check_out_time"`
	CheckInLat      *float64   `gorm:"type:decimal(10,8)" json:"check_in_lat"`
	CheckInLng      *float64   `gorm:"type:decimal(11,8)" json:"check_in_lng"`
	CheckOutLat     *float64   `gorm:"type:decimal(10,8)" json:"check_out_lat"`
	CheckOutLng     *float64   `gorm:"type:decimal(11,8)" json:"check_out_lng"`
	CheckInAddress  string     `gorm:"type:text" json:"check_in_address"`
	CheckOutAddress string     `gorm:"type:text" json:"check_out_address"`
	Status          string     `gorm:"type:varchar(20);default:hadir" json:"status"`
	Notes           string     `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (p *Presensi) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
