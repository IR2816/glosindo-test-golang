package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID          string     `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID      string     `gorm:"type:varchar(50);not null;index" json:"user_id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Category    string     `gorm:"type:varchar(50);not null;index" json:"category"`
	Priority    string     `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Status      string     `gorm:"type:varchar(20);default:'open';index" json:"status"`
	AssignedTo  *string    `gorm:"type:varchar(50)" json:"assigned_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at"`

	User     User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Progress []TicketProgress `gorm:"foreignKey:TicketID" json:"progress,omitempty"`
}

type TicketProgress struct {
	ID        string    `gorm:"type:varchar(50);primaryKey" json:"id"`
	TicketID  string    `gorm:"type:varchar(50);not null;index" json:"ticket_id"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedBy string    `gorm:"type:varchar(50);not null" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`

	Ticket Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"-"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (tp *TicketProgress) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == "" {
		tp.ID = uuid.New().String()
	}
	return nil
}