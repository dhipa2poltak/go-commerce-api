package domain

import (
	"time"
)

type Product struct {
	ID          int64     `gorm:primaryKey`
	CategoryID  int64     `gorm:"not null;index"`
	Name        string    `gorm:"size:200;not null"`
	Description string    `gorm:"type:text"`
	Price       int64     `gorm:"type:numeric(19,2);not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`

	Category Category `gorm:"foreignKey:categoryID"`
}
