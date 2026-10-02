APP      := api
DB_URL   ?= postgres://app:app@localhost:5432/app?sslmode=disable

.PHONY: run build test lint tidy up down migrate-up migrate-down

run:
	set -a; . ./.env; set +a; go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/api

test:
	go test -race -cover ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

up:
	docker compose -f deployments/docker-compose.yml up --build

down:
	docker compose -f deployments/docker-compose.yml down -v

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
