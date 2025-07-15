package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusCreated    OrderStatus = "created"
	StatusAssigned   OrderStatus = "assigned"
	StatusDelivering OrderStatus = "delivering"
	StatusCompleted  OrderStatus = "completed"
	StatusCanceled   OrderStatus = "canceled"
)

type Order struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	UserID          uuid.UUID   `json:"user_id" db:"user_id" validate:"required"`
	Status          OrderStatus `json:"status" db:"status" validate:"oneof=created assigned delivering completed canceled"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
	ETAMinutes      *int        `json:"eta_minutes,omitempty" db:"eta_minutes"`
	DeliveryAddress string      `json:"delivery_address" db:"delivery_address" validate:"required,min=10,max=500"`
}

// Для создания нового заказа (без ID и временных меток)
type CreateOrderRequest struct {
	UserID          uuid.UUID   `json:"user_id" validate:"required"`
	DeliveryAddress string      `json:"delivery_address" validate:"required,min=10,max=500"`
	Status          OrderStatus `json:"status" validate:"omitempty,oneof=created assigned delivering completed canceled"`
}

// Для обновления заказа
type UpdateOrderRequest struct {
	Status     OrderStatus `json:"status" validate:"omitempty,oneof=assigned delivering completed canceled"`
	ETAMinutes *int        `json:"eta_minutes" validate:"omitempty,min=1,max=1440"` // 1 мин - 24 часа
}
