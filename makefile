SWAGGER_OUT=./docs

POSTGRES_USER=postgres
POSTGRES_PASSWORD=1234
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=auth_service_db

# DB_URL=postgres://user:password@host:port/db?sslmode=disable
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable


.PHONY: swag run migrate up down create_migration

swag:
	swag init --generalInfo ./cmd/server/main.go --output $(SWAGGER_OUT)

run:
	go run ./cmd/server/main.go

migrate:
	migrate -path ./migrations -database "$(DB_URL)" up

up:
	migrate -path ./migrations -database "$(DB_URL)" up

down:
	migrate -path ./migrations -database "$(DB_URL)" down

create_migration:
	migrate create -ext sql -dir ./migrations -seq init_schema