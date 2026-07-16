package domain

import "time"

// ItemIngredientRef is a lightweight reference to an ingredient used in item payloads.
type ItemIngredientRef struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Item is the domain entity for tm_item.
type Item struct {
	UUID        string              `json:"uuid"`
	Name        string              `json:"name"`
	Price       float64             `json:"price"`
	Status      RecordStatus        `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty"`
	DeletedAt   *time.Time          `json:"deleted_at,omitempty"`
	Ingredients []ItemIngredientRef `json:"ingredients,omitempty"`
}

// ItemCreateInput is the payload for creating an item with its ingredients.
type ItemCreateInput struct {
	Name        string       `json:"name"`
	Price       float64      `json:"price"`
	Status      RecordStatus `json:"status"`
	Ingredients []string     `json:"ingredients"`
}

// ItemUpdateInput is the payload for updating an item and its ingredients.
type ItemUpdateInput struct {
	UUID        string       `json:"-"`
	Name        string       `json:"name"`
	Price       float64      `json:"price"`
	Status      RecordStatus `json:"status"`
	Ingredients []string     `json:"ingredients"`
}

// PaginatedIngredients holds a page of ingredient results.
type PaginatedIngredients struct {
	Data       []Ingredient `json:"data"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	Total      int          `json:"total"`
}

// PaginatedItems holds a page of item results.
type PaginatedItems struct {
	Data     []Item `json:"data"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Total    int    `json:"total"`
}
