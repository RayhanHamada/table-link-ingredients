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

// ItemRepository defines the data-access contract for tm_item.
type ItemRepository interface {
	FindAll(ctx context.Context) ([]domain.Item, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Item, error)
	Create(ctx context.Context, item *domain.Item) error
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, uuid string) error
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type itemRepo struct {
	pool *pgxpool.Pool
}

// NewItemRepository returns the concrete implementation.
func NewItemRepository(pool *pgxpool.Pool) ItemRepository {
	return &itemRepo{pool: pool}
}

// ---------------------------------------------------------------------------
// Queries
// ---------------------------------------------------------------------------

const (
	itemFindAllSQL = `
		SELECT uuid, name, price, status, created_at, updated_at, deleted_at
		FROM tm_item
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	itemFindByUUIDSQL = `
		SELECT uuid, name, price, status, created_at, updated_at, deleted_at
		FROM tm_item
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	itemCreateSQL = `
		INSERT INTO tm_item (uuid, name, price, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	itemUpdateSQL = `
		UPDATE tm_item
		SET name = $1, price = $2, status = $3, updated_at = NOW()
		WHERE uuid = $4 AND deleted_at IS NULL
	`

	itemDeleteSQL = `
		UPDATE tm_item SET deleted_at = NOW() WHERE uuid = $1
	`
)

// ---------------------------------------------------------------------------
// pgx.Rows scanning helpers
// ---------------------------------------------------------------------------

func scanItem(row pgx.Row) (*domain.Item, error) {
	var i domain.Item
	err := row.Scan(
		&i.UUID,
		&i.Name,
		&i.Price,
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

func scanItems(rows pgx.Rows) ([]domain.Item, error) {
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var i domain.Item
		if err := rows.Scan(
			&i.UUID,
			&i.Name,
			&i.Price,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Method implementations
// ---------------------------------------------------------------------------

func (r *itemRepo) FindAll(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx, itemFindAllSQL)
	if err != nil {
		return nil, err
	}
	return scanItems(rows)
}

func (r *itemRepo) FindByUUID(ctx context.Context, uuid string) (*domain.Item, error) {
	row := r.pool.QueryRow(ctx, itemFindByUUIDSQL, uuid)
	return scanItem(row)
}

func (r *itemRepo) Create(ctx context.Context, item *domain.Item) error {
	_, err := r.pool.Exec(ctx, itemCreateSQL, item.UUID, item.Name, item.Price, item.Status)
	return err
}

func (r *itemRepo) Update(ctx context.Context, item *domain.Item) error {
	_, err := r.pool.Exec(ctx, itemUpdateSQL, item.Name, item.Price, item.Status, item.UUID)
	return err
}

func (r *itemRepo) Delete(ctx context.Context, uuid string) error {
	_, err := r.pool.Exec(ctx, itemDeleteSQL, uuid)
	return err
}
