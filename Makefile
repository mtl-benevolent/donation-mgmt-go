CMD_PATH = "src/cmd"
DIST_PATH = "dist"

.PHONY: build
build: generate
	@echo "Building the API"
	@go build -o $(DIST_PATH)/api $(CMD_PATH)/api.go

.PHONY: build_debug
build_debug:
	@echo "Building the API in debug mode"
	@go build -gcflags="all=-N -l" -o $(DIST_PATH)/api $(CMD_PATH)/api.go

generate:
	@echo "[INFO] Generating code..."
	@go generate ./...

schema.gen.sql: ./prisma/schema.prisma
	@echo "[INFO] Generating SQL Schema from Prisma"
	@npx prisma migrate diff --script --from-empty --to-schema-datamodel=./prisma/schema.prisma > schema.gen.sql

.PHONY: sqlc
sqlc:
	@echo "[INFO] Installing sqlc"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.23.0

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
	@echo "Installing linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	
	@echo "Linting code"
	@$$(go env GOPATH)/bin/golangci-lint run ./...
