package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginHistory struct {
	ID         string     `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID     string     `gorm:"type:varchar(50);not null;index" json:"user_id"`
	LoginTime  time.Time  `gorm:"not null" json:"login_time"`
	LogoutTime *time.Time `json:"logout_time"`
	DeviceInfo string     `gorm:"type:text" json:"device_info"`
	IPAddress  string     `gorm:"type:varchar(50)" json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
	
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (l *LoginHistory) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}