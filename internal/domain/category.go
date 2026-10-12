package domain

import (
	"time"
)

type Category struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
