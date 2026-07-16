package repository

import (
	"context"

	"tablelink-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// ItemIngredientRepository defines the data-access contract for
// the join table tm_item_ingredient. Relationships use hard delete.
type ItemIngredientRepository interface {
	FindByItemUUID(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error)
	FindRefsByItemUUID(ctx context.Context, itemUUID string) ([]domain.ItemIngredientRef, error)
	CreateBulkTx(ctx context.Context, tx pgx.Tx, itemUUID string, ingredientUUIDs []string) error
	DeleteByItemUUIDTx(ctx context.Context, tx pgx.Tx, itemUUID string) error
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

	itemIngredientCreateSQL = `
		INSERT INTO tm_item_ingredient (uuid_item, uuid_ingredient) VALUES ($1, $2)
	`

	itemIngredientDeleteByItemSQL = `
		DELETE FROM tm_item_ingredient WHERE uuid_item = $1
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

// FindRefsByItemUUID returns label/value refs by joining with tm_ingredient.
func (r *itemIngredientRepo) FindRefsByItemUUID(ctx context.Context, itemUUID string) ([]domain.ItemIngredientRef, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT i.name, ii.uuid_ingredient
		 FROM tm_item_ingredient ii
		 JOIN tm_ingredient i ON i.uuid = ii.uuid_ingredient
		 WHERE ii.uuid_item = $1 AND i.deleted_at IS NULL`, itemUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []domain.ItemIngredientRef
	for rows.Next() {
		var ref domain.ItemIngredientRef
		if err := rows.Scan(&ref.Label, &ref.Value); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, rows.Err()
}

// CreateBulkTx inserts all ingredient relationships for an item inside a
// transaction. Uses hard insert.
func (r *itemIngredientRepo) CreateBulkTx(ctx context.Context, tx pgx.Tx, itemUUID string, ingredientUUIDs []string) error {
	for _, ingUUID := range ingredientUUIDs {
		if _, err := tx.Exec(ctx, itemIngredientCreateSQL, itemUUID, ingUUID); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByItemUUIDTx performs a hard delete of all relationships for a given
// item inside a transaction.
func (r *itemIngredientRepo) DeleteByItemUUIDTx(ctx context.Context, tx pgx.Tx, itemUUID string) error {
	_, err := tx.Exec(ctx, itemIngredientDeleteByItemSQL, itemUUID)
	return err
}

