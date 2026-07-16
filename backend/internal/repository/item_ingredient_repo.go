package repository

import (
	"context"

	"tablelink-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// ItemIngredientRepository defines the data-access contract for
// the join table tm_item_ingredient.
type ItemIngredientRepository interface {
	FindByItemUUID(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error)
	FindByIngredientUUID(ctx context.Context, ingredientUUID string) ([]domain.ItemIngredient, error)
	Create(ctx context.Context, rel *domain.ItemIngredient) error
	Delete(ctx context.Context, itemUUID, ingredientUUID string) error
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type itemIngredientRepo struct {
	pool *pgxpool.Pool
}

// NewItemIngredientRepository returns the concrete implementation.
func NewItemIngredientRepository(pool *pgxpool.Pool) ItemIngredientRepository {
	return &itemIngredientRepo{pool: pool}
}

// ---------------------------------------------------------------------------
// Queries
// ---------------------------------------------------------------------------

const (
	itemIngredientFindByItemSQL = `
		SELECT uuid_item, uuid_ingredient
		FROM tm_item_ingredient
		WHERE uuid_item = $1
	`

	itemIngredientFindByIngredientSQL = `
		SELECT uuid_item, uuid_ingredient
		FROM tm_item_ingredient
		WHERE uuid_ingredient = $1
	`

	itemIngredientCreateSQL = `
		INSERT INTO tm_item_ingredient (uuid_item, uuid_ingredient) VALUES ($1, $2)
	`

	itemIngredientDeleteSQL = `
		DELETE FROM tm_item_ingredient WHERE uuid_item = $1 AND uuid_ingredient = $2
	`
)

// ---------------------------------------------------------------------------
// Method implementations
// ---------------------------------------------------------------------------

func (r *itemIngredientRepo) FindByItemUUID(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error) {
	rows, err := r.pool.Query(ctx, itemIngredientFindByItemSQL, itemUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []domain.ItemIngredient
	for rows.Next() {
		var rel domain.ItemIngredient
		if err := rows.Scan(&rel.UUIDItem, &rel.UUIDIngredient); err != nil {
			return nil, err
		}
		rels = append(rels, rel)
	}
	return rels, rows.Err()
}

func (r *itemIngredientRepo) FindByIngredientUUID(ctx context.Context, ingredientUUID string) ([]domain.ItemIngredient, error) {
	rows, err := r.pool.Query(ctx, itemIngredientFindByIngredientSQL, ingredientUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []domain.ItemIngredient
	for rows.Next() {
		var rel domain.ItemIngredient
		if err := rows.Scan(&rel.UUIDItem, &rel.UUIDIngredient); err != nil {
			return nil, err
		}
		rels = append(rels, rel)
	}
	return rels, rows.Err()
}

func (r *itemIngredientRepo) Create(ctx context.Context, rel *domain.ItemIngredient) error {
	_, err := r.pool.Exec(ctx, itemIngredientCreateSQL, rel.UUIDItem, rel.UUIDIngredient)
	return err
}

func (r *itemIngredientRepo) Delete(ctx context.Context, itemUUID, ingredientUUID string) error {
	_, err := r.pool.Exec(ctx, itemIngredientDeleteSQL, itemUUID, ingredientUUID)
	return err
}
