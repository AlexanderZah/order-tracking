package event

import "github.com/google/uuid"

type OrderCreated struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
	Address string    `json:"address"`
}

type OrderETAUpdated struct {
	OrderID    uuid.UUID `json:"order_id"`
	ETAMinutes int       `json:"eta_minutes"`
}
