API_ENTRY = "src/cmd/api/api.go"
DIST_PATH = "dist"

.PHONY: build
build: generate
	@echo "Building the API"
	@go build -o $(DIST_PATH)/api $(API_ENTRY)

.PHONY: build_debug
build_debug:
	@echo "Building the API in debug mode"
	@go build -gcflags="all=-N -l" -o $(DIST_PATH)/api $(API_ENTRY)

.PHONY: watch
watch:
	@echo "Starting the project in watch mode"
	@docker compose watch

.PHONE: cleanup
cleanup:
	@echo "Cleaning up local environment"
	@docker compose down -v

generate:
	@echo "[INFO] Generating code..."
	@go generate ./...

schema.gen.sql: ./prisma/schema.prisma
	@echo "[INFO] Generating SQL Schema from Prisma"
	@npx prisma migrate diff --script --from-empty --to-schema-datamodel=./prisma/schema.prisma > schema.gen.sql

.PHONY: sqlc
sqlc: schema.gen.sql
	@echo "Generating Data Access Layer using sqlc"
	@$$(go env GOPATH)/bin/sqlc generate

.PHONY: clean
clean:
	@echo "Cleaning up the build artifacts"
	@rm -rf dist
	@rm -rf **/*.gen.go
	@rm -rf **/*.gen.sql

.PHONY: lint
lint:
	@echo "Linting code"
	@$$(go env GOPATH)/bin/golangci-lint run ./...

.PHONY: seed
seed:
	@echo "Seeding the database"
	@go run ./src/cmd/seed/seed.go

.PHONY: deps
deps:
	@echo "[INFO] Installing linter"
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

	@echo "[INFO] Installing sqlc"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0

.PHONY: migrate
migrate:
	@echo "[INFO] Applying pending migrations (in dev mode)"
	@npx prisma migrate dev
