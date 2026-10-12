package domain

import (
	"time"
)

type OrderItem struct {
	ID        int64     `gorm:"primaryKey"`
	OrderID   int64     `gorm:"not null;index"`
	ProductID int64     `gorm:"not null;index"`
	Quantity  int       `gorm:"not null"`
	UnitPrice int64     `gorm:"type:numeric(19,2);not null"`
	Subtotal  int64     `gorm:"type:numeric(19,2);not null"`
	CreatedAt time.Time `gorm:"not null"`

	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
