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

// IngredientRepository defines the data-access contract for tm_ingredient.
type IngredientRepository interface {
	FindAll(ctx context.Context) ([]domain.Ingredient, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Ingredient, error)
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	Update(ctx context.Context, ingredient *domain.Ingredient) error
	Delete(ctx context.Context, uuid string) error
}

// ---------------------------------------------------------------------------
// Implementation (using pgxpool + pgx.Rows)
// ---------------------------------------------------------------------------

type ingredientRepo struct {
	pool *pgxpool.Pool
}

// NewIngredientRepository returns the concrete implementation.
func NewIngredientRepository(pool *pgxpool.Pool) IngredientRepository {
	return &ingredientRepo{pool: pool}
}

// ---------------------------------------------------------------------------
// Queries
// ---------------------------------------------------------------------------

const (
	ingredientFindAllSQL = `
		SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at
		FROM tm_ingredient
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	ingredientFindByUUIDSQL = `
		SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at
		FROM tm_ingredient
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	ingredientCreateSQL = `
		INSERT INTO tm_ingredient (uuid, name, cause_alergy, type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	ingredientUpdateSQL = `
		UPDATE tm_ingredient
		SET name = $1, cause_alergy = $2, type = $3, status = $4, updated_at = NOW()
		WHERE uuid = $5 AND deleted_at IS NULL
	`

	ingredientDeleteSQL = `
		UPDATE tm_ingredient SET deleted_at = NOW() WHERE uuid = $1
	`
)

// ---------------------------------------------------------------------------
// pgx.Rows scanning helpers
// ---------------------------------------------------------------------------

// scanIngredient scans a single pgx.Row into a domain.Ingredient.
func scanIngredient(row pgx.Row) (*domain.Ingredient, error) {
	var i domain.Ingredient
	err := row.Scan(
		&i.UUID,
		&i.Name,
		&i.CauseAlergy,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// scanIngredients iterates pgx.Rows and collects all scanned ingredients.
func scanIngredients(rows pgx.Rows) ([]domain.Ingredient, error) {
	defer rows.Close()

	var ingredients []domain.Ingredient
	for rows.Next() {
		var i domain.Ingredient
		if err := rows.Scan(
			&i.UUID,
			&i.Name,
			&i.CauseAlergy,
			&i.Type,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, rows.Err()
}

// ---------------------------------------------------------------------------
// Method implementations
// ---------------------------------------------------------------------------

func (r *ingredientRepo) FindAll(ctx context.Context) ([]domain.Ingredient, error) {
	rows, err := r.pool.Query(ctx, ingredientFindAllSQL)
	if err != nil {
		return nil, err
	}
	return scanIngredients(rows)
}

func (r *ingredientRepo) FindByUUID(ctx context.Context, uuid string) (*domain.Ingredient, error) {
	row := r.pool.QueryRow(ctx, ingredientFindByUUIDSQL, uuid)
	return scanIngredient(row)
}

func (r *ingredientRepo) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	_, err := r.pool.Exec(ctx, ingredientCreateSQL,
		ingredient.UUID,
		ingredient.Name,
		ingredient.CauseAlergy,
		ingredient.Type,
		ingredient.Status,
	)
	return err
}

func (r *ingredientRepo) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	_, err := r.pool.Exec(ctx, ingredientUpdateSQL,
		ingredient.Name,
		ingredient.CauseAlergy,
		ingredient.Type,
		ingredient.Status,
		ingredient.UUID,
	)
	return err
}

func (r *ingredientRepo) Delete(ctx context.Context, uuid string) error {
	_, err := r.pool.Exec(ctx, ingredientDeleteSQL, uuid)
	return err
}
