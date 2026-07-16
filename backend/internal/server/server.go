package server

import (
	"context"
	"fmt"
	"log"

	"tablelink-backend/internal/config"
	"tablelink-backend/internal/handler"
	"tablelink-backend/internal/repository"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server is the application runtime. It owns the configuration, dependencies,
// and the Fiber app instance.
type Server struct {
	cfg   *config.Config
	app   *fiber.App
	pool  *pgxpool.Pool
	store *repository.Store
}

// New creates and wires every layer: config → db pool → repos → usecases → handlers → routes.
func New(cfg *config.Config) (*Server, error) {
	// ------------------------------------------------------------------
	// PostgreSQL connection pool (pgxpool)
	// ------------------------------------------------------------------
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	poolCfg.MaxConns = int32(cfg.Database.PoolMax)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify connectivity.
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Printf("connected to PostgreSQL at %s:%d/%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// ------------------------------------------------------------------
	// Dependency injection – wire up the layers bottom-to-top
	//   pool → store (repos) → usecases → handlers
	// ------------------------------------------------------------------
	store := repository.NewStore(pool)

	// Usecases
	ingredientUC := usecase.NewIngredientUsecase(store.Ingredient)
	itemUC := usecase.NewItemUsecase(store.Pool, store.Item, store.Ingredient, store.ItemIngredient)
	itemIngredientUC := usecase.NewItemIngredientUsecase(store.ItemIngredient)

	// Handlers
	ingredientH := handler.NewIngredientHandler(ingredientUC)
	itemH := handler.NewItemHandler(itemUC)
	itemIngredientH := handler.NewItemIngredientHandler(itemIngredientUC)

	// ------------------------------------------------------------------
	// Fiber app
	// ------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: defaultErrorHandler,
	})

	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger
	app.Use("/docs", static.New("./docs"))
	app.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
		httpSwagger.DeepLinking(true),
	)))

	// API v1 group
	v1 := app.Group("/api/v1")

	// Register route groups from each handler
	ingredientH.Register(v1)
	itemH.Register(v1)
	itemIngredientH.Register(itemH.RegisterGroup(v1)) // nested: /items/:uuid/ingredients

	return &Server{
		cfg:   cfg,
		app:   app,
		pool:  pool,
		store: store,
	}, nil
}

// Run starts the HTTP listener and blocks.
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	return s.app.Listen(addr)
}

// Shutdown gracefully stops the server and closes the database pool.
func (s *Server) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		log.Printf("app shutdown error: %v", err)
	}
	s.pool.Close()
	log.Println("server stopped")
	return nil
}

// ---------------------------------------------------------------------------
// Fiber default error handler
// ---------------------------------------------------------------------------

func defaultErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
