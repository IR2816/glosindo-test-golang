package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:varchar(50);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Address   string    `gorm:"type:text" json:"address"`
	Division  string    `gorm:"type:varchar(50)" json:"division"`
	Position  string    `gorm:"type:varchar(50)" json:"position"`
	JoinDate  *time.Time `gorm:"type:date" json:"join_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}