package notification

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	Message   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:varchar(10);default:'new'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
