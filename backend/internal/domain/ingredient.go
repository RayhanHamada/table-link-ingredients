package domain

import "time"

// IngredientType represents the dietary category of an ingredient.
// 0 = none, 1 = veggie, 2 = vegan.
type IngredientType int

const (
	IngredientTypeNone   IngredientType = 0
	IngredientTypeVeggie IngredientType = 1
	IngredientTypeVegan  IngredientType = 2
)

// RecordStatus is a soft-delete status.
// 0 = inactive, 1 = active.
type RecordStatus int

const (
	StatusInactive RecordStatus = 0
	StatusActive   RecordStatus = 1
)

// Ingredient is the domain entity for tm_ingredient.
type Ingredient struct {
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`
	CauseAlergy bool           `json:"cause_alergy"`
	Type        IngredientType `json:"type"`
	Status      RecordStatus   `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}
