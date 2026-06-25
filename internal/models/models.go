package models

import "time"

type Product struct {
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	CurrentStock int    `json:"current_stock"`
}

type Event struct {
	EventID    string    `json:"event_id"`
	SKU        string    `json:"sku"`
	Type       string    `json:"type"` // "IN" or "OUT"
	Quantity   int       `json:"quantity"`
	OccurredAt time.Time `json:"occurred_at"`
}
