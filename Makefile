CMD_PATH = "src/cmd"
DIST_PATH = "dist"

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

