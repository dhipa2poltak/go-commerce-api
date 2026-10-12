package domain

import (
	"time"
)

type Inventory struct {
	ID        int64     `gorm:"primaryKey"`
	ProductID int64     `gorm:"not null;uniqueIndex"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Product Product `gorm:"foreignKey:ProductID"`
}

func (Inventory) TableName() string {
	return "inventory"
}
