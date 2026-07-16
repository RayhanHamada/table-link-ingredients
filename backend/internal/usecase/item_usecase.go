package usecase

import (
	"context"
	"errors"
	"fmt"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemUsecase interface {
	List(ctx context.Context, page, pageSize int) (*domain.PaginatedItems, error)
	Get(ctx context.Context, uuid string) (*domain.Item, error)
	Create(ctx context.Context, input domain.ItemCreateInput) (*domain.Item, error)
	Update(ctx context.Context, input domain.ItemUpdateInput) (*domain.Item, error)
	Delete(ctx context.Context, uuid string) error
}

type itemUC struct {
	pool               *pgxpool.Pool
	itemRepo           repository.ItemRepository
	ingredientRepo     repository.IngredientRepository
	itemIngredientRepo repository.ItemIngredientRepository
}

func NewItemUsecase(
	pool *pgxpool.Pool,
	itemRepo repository.ItemRepository,
	ingredientRepo repository.IngredientRepository,
	itemIngredientRepo repository.ItemIngredientRepository,
) ItemUsecase {
	return &itemUC{
		pool: pool, itemRepo: itemRepo,
		ingredientRepo: ingredientRepo, itemIngredientRepo: itemIngredientRepo,
	}
}

func (uc *itemUC) List(ctx context.Context, page, pageSize int) (*domain.PaginatedItems, error) {
	data, total, err := uc.itemRepo.FindAllPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &domain.PaginatedItems{
		Data: data, Page: page, PageSize: pageSize, Total: total,
	}, nil
}

func (uc *itemUC) Get(ctx context.Context, uuid string) (*domain.Item, error) {
	item, err := uc.itemRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	refs, err := uc.itemIngredientRepo.FindRefsByItemUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	item.Ingredients = refs
	return item, nil
}

// --- Create (tx) ---
func (uc *itemUC) Create(ctx context.Context, input domain.ItemCreateInput) (*domain.Item, error) {
	if err := uc.validateName(ctx, input.Name, ""); err != nil {
		return nil, err
	}
	if len(input.Ingredients) == 0 {
		return nil, fmt.Errorf("at least one ingredient is required")
	}
	if _, err := uc.ingredientRepo.BatchExist(ctx, input.Ingredients); err != nil {
		return nil, fmt.Errorf("ingredient validation failed: %w", err)
	}
	item := &domain.Item{
		UUID: uuid.New().String(), Name: input.Name, Price: input.Price, Status: input.Status,
	}
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := uc.itemRepo.CreateTx(ctx, tx, item); err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}
	if err := uc.itemIngredientRepo.CreateBulkTx(ctx, tx, item.UUID, input.Ingredients); err != nil {
		return nil, fmt.Errorf("create relations: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	refs, err := uc.itemIngredientRepo.FindRefsByItemUUID(ctx, item.UUID)
	if err != nil {
		return nil, err
	}
	item.Ingredients = refs
	return item, nil
}

// --- Update (tx) ---
func (uc *itemUC) Update(ctx context.Context, input domain.ItemUpdateInput) (*domain.Item, error) {
	if err := uc.validateName(ctx, input.Name, input.UUID); err != nil {
		return nil, err
	}
	if len(input.Ingredients) == 0 {
		return nil, fmt.Errorf("at least one ingredient is required")
	}
	if _, err := uc.ingredientRepo.BatchExist(ctx, input.Ingredients); err != nil {
		return nil, fmt.Errorf("ingredient validation failed: %w", err)
	}
	item := &domain.Item{
		UUID: input.UUID, Name: input.Name, Price: input.Price, Status: input.Status,
	}
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := uc.itemRepo.UpdateTx(ctx, tx, item); err != nil {
		return nil, fmt.Errorf("update item: %w", err)
	}
	if err := uc.itemIngredientRepo.DeleteByItemUUIDTx(ctx, tx, item.UUID); err != nil {
		return nil, fmt.Errorf("delete old relations: %w", err)
	}
	if err := uc.itemIngredientRepo.CreateBulkTx(ctx, tx, item.UUID, input.Ingredients); err != nil {
		return nil, fmt.Errorf("create new relations: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	refs, err := uc.itemIngredientRepo.FindRefsByItemUUID(ctx, item.UUID)
	if err != nil {
		return nil, err
	}
	item.Ingredients = refs
	return item, nil
}

func (uc *itemUC) Delete(ctx context.Context, uuid string) error {
	return uc.itemRepo.Delete(ctx, uuid)
}

func (uc *itemUC) validateName(ctx context.Context, name, currentUUID string) error {
	existing, err := uc.itemRepo.FindByName(ctx, name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if existing != nil && existing.UUID != currentUUID {
		return fmt.Errorf("item name %q already exists", name)
	}
	return nil
}
