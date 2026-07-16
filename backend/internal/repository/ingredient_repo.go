package repository

import (
	"context"
	"fmt"

	"tablelink-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// IngredientRepository defines the data-access contract for tm_ingredient.
type IngredientRepository interface {
	FindAllPaginated(ctx context.Context, page, pageSize int) ([]domain.Ingredient, int, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Ingredient, error)
	FindByName(ctx context.Context, name string) (*domain.Ingredient, error)
	BatchExist(ctx context.Context, uuids []string) (bool, error)
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
	ingredientCountSQL = `
		SELECT COUNT(*) FROM tm_ingredient WHERE deleted_at IS NULL
	`

	ingredientFindAllPaginatedSQL = `
		SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at
		FROM tm_ingredient
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	ingredientFindByUUIDSQL = `
		SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at
		FROM tm_ingredient
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	ingredientFindByNameSQL = `
		SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at
		FROM tm_ingredient
		WHERE name = $1 AND deleted_at IS NULL
	`

	ingredientBatchExistSQL = `
		SELECT COUNT(*) FROM tm_ingredient
		WHERE uuid = ANY($1) AND deleted_at IS NULL
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

func scanIngredient(row pgx.Row) (*domain.Ingredient, error) {
	var i domain.Ingredient
	err := row.Scan(
		&i.UUID, &i.Name, &i.CauseAlergy, &i.Type, &i.Status,
		&i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func scanIngredients(rows pgx.Rows) ([]domain.Ingredient, error) {
	defer rows.Close()
	var ingredients []domain.Ingredient
	for rows.Next() {
		var i domain.Ingredient
		if err := rows.Scan(
			&i.UUID, &i.Name, &i.CauseAlergy, &i.Type, &i.Status,
			&i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
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

func (r *ingredientRepo) FindAllPaginated(ctx context.Context, page, pageSize int) ([]domain.Ingredient, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, ingredientCountSQL).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, ingredientFindAllPaginatedSQL, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	data, err := scanIngredients(rows)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *ingredientRepo) FindByUUID(ctx context.Context, uuid string) (*domain.Ingredient, error) {
	row := r.pool.QueryRow(ctx, ingredientFindByUUIDSQL, uuid)
	return scanIngredient(row)
}

func (r *ingredientRepo) FindByName(ctx context.Context, name string) (*domain.Ingredient, error) {
	row := r.pool.QueryRow(ctx, ingredientFindByNameSQL, name)
	return scanIngredient(row)
}

func (r *ingredientRepo) BatchExist(ctx context.Context, uuids []string) (bool, error) {
	if len(uuids) == 0 {
		return true, nil
	}
	var count int
	if err := r.pool.QueryRow(ctx, ingredientBatchExistSQL, uuids).Scan(&count); err != nil {
		return false, err
	}
	if count != len(uuids) {
		return false, fmt.Errorf("some ingredients do not exist: expected %d, found %d", len(uuids), count)
	}
	return true, nil
}

func (r *ingredientRepo) Create(ctx context.Context, i *domain.Ingredient) error {
	_, err := r.pool.Exec(ctx, ingredientCreateSQL,
		i.UUID, i.Name, i.CauseAlergy, i.Type, i.Status,
	)
	return err
}

func (r *ingredientRepo) Update(ctx context.Context, i *domain.Ingredient) error {
	_, err := r.pool.Exec(ctx, ingredientUpdateSQL,
		i.Name, i.CauseAlergy, i.Type, i.Status, i.UUID,
	)
	return err
}

func (r *ingredientRepo) Delete(ctx context.Context, uuid string) error {
	_, err := r.pool.Exec(ctx, ingredientDeleteSQL, uuid)
	return err
}

