package domain

import "time"

// Item is the domain entity for tm_item.
type Item struct {
	UUID      string       `json:"uuid"`
	Name      string       `json:"name"`
	Price     float64      `json:"price"`
	Status    RecordStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at,omitempty"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty"`
}
