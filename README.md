# TableLink

A full-stack application for managing ingredients and menu items with many-to-many relationships. Built with Go (Fiber) + React (TypeScript) following clean architecture principles.

## Architecture

```
tablelink/
├── backend/                  # Go REST API
│   ├── cmd/server/main.go    # Entry point
│   ├── internal/
│   │   ├── config/           # Environment configuration
│   │   ├── domain/           # Entities & value objects
│   │   ├── repository/       # Data access (pgx + PostgreSQL)
│   │   ├── usecase/          # Business logic
│   │   ├── handler/          # HTTP handlers (Fiber)
│   │   └── server/           # DI wiring & app assembly
│   ├── migrations/           # Goose SQL migrations
│   └── docker-compose.yaml   # PostgreSQL for development
│
└── frontend/                 # React + TypeScript SPA
    └── src/
        ├── services/         # API client
        ├── hooks/            # TanStack Query hooks
        ├── lib/              # Utilities & Valibot schemas
        ├── components/       # Reusable shared components
        └── features/         # Business feature modules (ingredients, items)
```

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | ≥ 1.23 | Backend runtime |
| Docker + Compose | latest | PostgreSQL for development |
| Goose | latest | Database migrations |
| Bun | latest | Frontend package manager (or npm/yarn) |
| Make | — | Task runner |

### Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Quick Start

### 1. Start PostgreSQL

```bash
cd backend
docker compose up -d
```

PostgreSQL 16 starts on port `5432` with database `tablelink`, user `postgres`, password `postgres`.

### 2. Run migrations

```bash
cd backend
make migrate-up
```

This creates the three tables (`tm_ingredient`, `tm_item`, `tm_item_ingredient`) with primary keys, foreign keys, and seed data.

### 3. Start the backend

```bash
cd backend
make run
```

The API starts at `http://localhost:3000`.

Swagger UI: `http://localhost:3000/swagger/index.html`

### 4. Start the frontend

```bash
cd frontend
bun install
bun run dev
```

The app starts at `http://localhost:5173`.

## Makefile Commands

### Backend

| Command | Description |
|---|---|
| `make build` | Build binary to `bin/app` |
| `make run` | Build and run with `go run` |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-reset` | Down then up (wipe + rebuild) |
| `make swagger` | Regenerate Swagger docs |

### Frontend

| Command | Description |
|---|---|
| `bun run dev` | Start dev server (port 5173) |
| `bun run build` | Production build |
| `bun run preview` | Preview production build |

## API Endpoints

Base URL: `http://localhost:3000/api/v1`

### Ingredients

| Method | Path | Description |
|---|---|---|
| `GET` | `/ingredients` | List (paginated: `?page=1&page_size=10`) |
| `GET` | `/ingredients/:uuid` | Get one |
| `POST` | `/ingredients` | Create |
| `PUT` | `/ingredients/:uuid` | Update |
| `DELETE` | `/ingredients/:uuid` | Soft delete |

### Items

| Method | Path | Description |
|---|---|---|
| `GET` | `/items` | List (paginated) |
| `GET` | `/items/:uuid` | Get one (includes ingredient UUIDs) |
| `POST` | `/items` | Create with ingredient relationships |
| `PUT` | `/items/:uuid` | Update and replace relationships |
| `DELETE` | `/items/:uuid` | Soft delete |

### Item Ingredients (read-only)

| Method | Path | Description |
|---|---|---|
| `GET` | `/items/:uuid/ingredients` | List ingredient UUIDs for an item |

## Database

### Tables

**tm_ingredient** — ingredient catalog

| Column | Type | Description |
|---|---|---|
| uuid | uuid | Primary key |
| name | varchar(255) | Unique name |
| cause_alergy | bool | Whether it causes allergies |
| type | int4 | 0 = None, 1 = Veggie, 2 = Vegan |
| status | int4 | 0 = Inactive, 1 = Active |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp | Soft delete |

**tm_item** — menu items

| Column | Type | Description |
|---|---|---|
| uuid | uuid | Primary key |
| name | varchar(255) | Unique name |
| price | numeric(10,2) | Price in IDR |
| status | int4 | 0 = Inactive, 1 = Active |
| created_at | timestamp | |
| updated_at | timestamp | |
| deleted_at | timestamp | Soft delete |

**tm_item_ingredient** — many-to-many junction

| Column | Type | Description |
|---|---|---|
| uuid_item | uuid | FK → tm_item |
| uuid_ingredient | uuid | FK → tm_ingredient |

### Seed Data

| Ingredient | Type |
|---|---|
| Chicken | None |
| Pork | None |
| Radish | Vegan |
| Egg | Veggie |

| Item | Price | Ingredients |
|---|---|---|
| Chicken Pork | 30,000 | Chicken, Pork |
| Chicken Pork with Radish | 35,000 | Chicken, Pork, Radish |
| Salad Egg | 20,000 | Radish, Egg |

## Environment Variables

Copy and edit `.env` in `backend/`:

| Variable | Default |
|---|---|
| `DB_HOST` | localhost |
| `DB_PORT` | 5432 |
| `DB_USER` | postgres |
| `DB_PASSWORD` | postgres |
| `DB_NAME` | tablelink |
| `DB_SCHEMA` | public |
| `DB_POOL_MAX` | 20 |
| `SERVER_PORT` | 3000 |
| `APP_NAME` | tablelink-backend |
| `APP_ENV` | development |
