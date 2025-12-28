package model

import (
	"time"

	"gorm.io/datatypes"
)

type Job struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Topic        string         `gorm:"index;size:255;not null" json:"topic"`
	Payload      datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	Status       string         `gorm:"index;size:50;default:'pending'" json:"status"` // pending, processing, completed, failed
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (Job) TableName() string {
	return "jobs"
}
