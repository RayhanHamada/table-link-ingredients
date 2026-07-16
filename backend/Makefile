build:
	go build -o bin/app ./cmd/server

run:
	go run ./cmd/server

DSN := $(shell set -a && . ./.env 2>/dev/null && echo "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable")

migrate-up:
	goose -dir migrations postgres "$(DSN)" up

migrate-down:
	goose -dir migrations postgres "$(DSN)" down

migrate-reset:
	goose -dir migrations postgres "$(DSN)" reset

swagger:
	swag init -g cmd/server/main.go -o docs
