# user-service

User registration and lookup service. Version: **v0.1.0**.

Contract-first REST API (OpenAPI), DDD with Clean Architecture layers, PostgreSQL.

## Stack

| Area | Technology |
|---|---|
| Language | Go 1.27 |
| API contract | OpenAPI 3.0.3, code generated with oapi-codegen v2.4.0 |
| HTTP router | chi v5 |
| Database | PostgreSQL 17 (database/sql + pgx v5) |
| Migrations | golang-migrate v4 (own migrator binary) |
| Config | app.ini (gopkg.in/ini.v1) + DATABASE_URL env variable |
| Tests | testify |
| Lint | golangci-lint v2 |
| API docs | Swagger UI |
| Infrastructure | Docker, Docker Compose |

## Requirements

- Go 1.27+
- Docker with Compose v2
- make

## Quick start

```
make up
```

This starts PostgreSQL, applies migrations (migrator) and runs the API.

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger
- OpenAPI contract: http://localhost:8080/openapi.yaml

Stop: `make down` (data is kept). Full reset: `docker compose down -v`.

## Local development

```
make init                          # create config/app.ini from the template
docker compose up -d postgres
make migrate-up
make run
```

## Commands

| Command | Description |
|---|---|
| `make generate` | Generate Go code from `api/openapi.yaml` |
| `make build` | Build all packages |
| `make test` | Run unit tests |
| `make lint` | Run golangci-lint |
| `make migrate-up` | Apply all new migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-version` | Show the current schema version |
| `make migrate-new name=add_roles` | Create a new migration pair |
| `make up` / `make down` | Start / stop everything in Docker Compose |

## Architecture

Layers (dependencies point inwards):

```
transport -> application -> domain <- infrastructure
```

- `domain`: User Aggregate Root, Email and PasswordHash Value Objects, domain rules and repository interface
- `application/service/user`: application service, use cases, DTO and ports
- `infrastructure`: PostgreSQL repository and bcrypt hasher implementations
- `transport/http`: HTTP handlers, mapper, router and generated API types

## Rules

- The API contract `api/openapi.yaml` is the source of truth. Change it first, then run `make generate`. Never edit `api.gen.go` by hand.
- The database schema is changed only through migrations (`migrations/`), applied only by the migrator.
- `config/app.ini` is not committed. Use `config/app.ini.example` as a template. Secrets go through environment variables.