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
