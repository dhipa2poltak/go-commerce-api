package domain

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
)

type Order struct {
	ID          int64       `gorm:"primaryKey"`
	UserID      int64       `gorm:"not null;index"`
	Status      OrderStatus `gorm:"size:30;not null"`
	TotalAmount int64       `gorm:"type:numeric(19,2);not null"`
	CreatedAt   time.Time   `gorm:"not null"`
	UpdatedAt   time.Time   `gorm:"not null"`

	User  User        `gorm:"foreignKey:UserID"`
	Items []OrderItem `gorm:"foreignKey:OrderID"`
}
