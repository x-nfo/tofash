package model

import "time"

type Notification struct {
	ID               uint    `gorm:"primaryKey"`
	ReceiverID       uint    `gorm:"status:not null"`
	Subject          *string `gorm:"type:varchar(255)"`
	Message          string  `gorm:"type:text"`
	Status           string  `gorm:"type:varchar(50)"` // SENT, PENDING, READ
	NotificationType string  `gorm:"type:varchar(50)"` // EMAIL, PUSH
	SentAt           *time.Time
	ReadAt           *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (Notification) TableName() string {
	return "notifications"
}
