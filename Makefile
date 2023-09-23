CMD_PATH = "src/cmd"
DIST_PATH = "dist"

DB_HOST := "${DB_HOST}"
ifeq ($(DB_HOST), "")
	DB_HOST := "localhost"
endif

DB_PORT := "${DB_PORT}"
ifeq ($(DB_PORT), "")
	DB_PORT := "26257"
endif

DB_MIGRATE_USER := "${DB_MIGRATE_USER}"
ifeq ($(DB_MIGRATE_USER), "")
	DB_MIGRATE_USER := "donation_mgmt_migrator"
endif

DB_MIGRATE_PASSWORD := "${DB_MIGRATE_PASSWORD}"

DB_NAME := "${DB_NAME}"
ifeq (${DB_NAME}, "")
	DB_NAME := "donationsdb"
endif

DB_SCHEMA := "${DB_SCHEMA}"
ifeq (${DB_SCHEMA}, "")
	DB_SCHEMA := "donations"
endif

DB_STRING = "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_MIGRATE_USER) password=$(DB_MIGRATE_PASSWORD) dbname=${DB_NAME} sslmode=disable"

.PHONY: build
build:
	@echo "Building the API"
	@go build -o $(DIST_PATH)/api $(CMD_PATH)/api.go


.PHONY: clean
clean:
	@echo "Cleaning up the build artifacts"
	@rm -rf dist

.PHONY: lint
lint:
	@echo "Installing linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	
	@echo "Linting code"
	@$$(go env GOPATH)/bin/golangci-lint run ./...

.PHONY: goose
goose:
	@echo "Installing Goose"
	@go install github.com/pressly/goose/v3/cmd/goose@v3.15.0

.PHONY: db_up
db_up: goose
	@echo "Applying migrations"
	@$$(go env GOPATH)/bin/goose -dir migrations -table "$(DB_NAME).$(DB_SCHEMA).goose_migrations" postgres $(DB_STRING) up

.PHONY: db_down
db_down: goose
	@echo "Reverting last migration"
	@$$(go env GOPATH)/bin/goose -dir migrations -table "$(DB_NAME).$(DB_SCHEMA).goose_migrations" postgres $(DB_STRING) down

.PHONY: db_status
db_status: goose
	@echo "Querying migration status"
	@$$(go env GOPATH)/bin/goose -dir migrations -table "$(DB_NAME).$(DB_SCHEMA).goose_migrations" postgres $(DB_STRING) status
