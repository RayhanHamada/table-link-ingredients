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
	FindAllPaginated(ctx context.Context, page, pageSize int) ([]domain.Item, int, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.Item, error)
	FindByName(ctx context.Context, name string) (*domain.Item, error)
	CreateTx(ctx context.Context, tx pgx.Tx, item *domain.Item) error
	UpdateTx(ctx context.Context, tx pgx.Tx, item *domain.Item) error
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
	itemCountSQL = `SELECT COUNT(*) FROM tm_item WHERE deleted_at IS NULL`

	itemFindAllPaginatedSQL = `
		SELECT uuid, name, price, status, created_at, updated_at, deleted_at
		FROM tm_item
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	itemFindByUUIDSQL = `
		SELECT uuid, name, price, status, created_at, updated_at, deleted_at
		FROM tm_item
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	itemFindByNameSQL = `
		SELECT uuid, name, price, status, created_at, updated_at, deleted_at
		FROM tm_item
		WHERE name = $1 AND deleted_at IS NULL
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

	itemDeleteSQL = `UPDATE tm_item SET deleted_at = NOW() WHERE uuid = $1`
)

// ---------------------------------------------------------------------------
// pgx.Rows scanning helpers
// ---------------------------------------------------------------------------

func scanItem(row pgx.Row) (*domain.Item, error) {
	var i domain.Item
	err := row.Scan(&i.UUID, &i.Name, &i.Price, &i.Status, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt)
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
		if err := rows.Scan(&i.UUID, &i.Name, &i.Price, &i.Status, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Method implementations
// ---------------------------------------------------------------------------

func (r *itemRepo) FindAllPaginated(ctx context.Context, page, pageSize int) ([]domain.Item, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, itemCountSQL).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, itemFindAllPaginatedSQL, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	data, err := scanItems(rows)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *itemRepo) FindByUUID(ctx context.Context, uuid string) (*domain.Item, error) {
	row := r.pool.QueryRow(ctx, itemFindByUUIDSQL, uuid)
	return scanItem(row)
}

func (r *itemRepo) FindByName(ctx context.Context, name string) (*domain.Item, error) {
	row := r.pool.QueryRow(ctx, itemFindByNameSQL, name)
	return scanItem(row)
}

func (r *itemRepo) CreateTx(ctx context.Context, tx pgx.Tx, i *domain.Item) error {
	_, err := tx.Exec(ctx, itemCreateSQL, i.UUID, i.Name, i.Price, i.Status)
	return err
}

func (r *itemRepo) UpdateTx(ctx context.Context, tx pgx.Tx, i *domain.Item) error {
	_, err := tx.Exec(ctx, itemUpdateSQL, i.Name, i.Price, i.Status, i.UUID)
	return err
}

func (r *itemRepo) Delete(ctx context.Context, uuid string) error {
	_, err := r.pool.Exec(ctx, itemDeleteSQL, uuid)
	return err
}

