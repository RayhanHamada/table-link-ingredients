package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store holds the database connection pool and all repository implementations.
// It is the single point of access to all data layers.
type Store struct {
	Ingredient      IngredientRepository
	Item            ItemRepository
	ItemIngredient  ItemIngredientRepository
}

// NewStore wires all repositories together with the provided pgxpool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Ingredient:      NewIngredientRepository(pool),
		Item:            NewItemRepository(pool),
		ItemIngredient:  NewItemIngredientRepository(pool),
	}
}
