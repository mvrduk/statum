.PHONY: init generate build test lint run migrate-up migrate-down migrate-version migrate-new up down logs

SERVICE_DIR := services/user-service

init:
	cp -n $(SERVICE_DIR)/config/app.ini.example $(SERVICE_DIR)/config/app.ini

generate:
	cd $(SERVICE_DIR) && go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

build:
	cd $(SERVICE_DIR) && go build ./...

test:
	cd $(SERVICE_DIR) && go test ./...

lint:
	cd $(SERVICE_DIR) && go tool golangci-lint run ./...

run:
	cd $(SERVICE_DIR) && go run ./cmd/api

migrate-up:
	cd $(SERVICE_DIR) && go run ./cmd/migrator up

migrate-down:
	cd $(SERVICE_DIR) && go run ./cmd/migrator down

migrate-version:
	cd $(SERVICE_DIR) && go run ./cmd/migrator version

migrate-new:
	@test -n "$(name)" || (echo "usage: make migrate-new name=add_roles" && exit 1)
	@ts=$$(date +%Y%m%d%H%M%S); \
	touch $(SERVICE_DIR)/migrations/$${ts}_$(name).up.sql $(SERVICE_DIR)/migrations/$${ts}_$(name).down.sql; \
	echo "created $${ts}_$(name)"

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f api