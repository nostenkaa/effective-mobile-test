# Name of the executable
APP_NAME=effective-mobile-test

# Directory for binaries
BIN_DIR=bin



# Show available commands
help: ## Show available commands
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove binaries"
	@echo "  docker   - Build and run Docker container"

build: ## Build the application
	@echo "Building application..."
ifeq ($(OS),Windows_NT)
	@if not exist $(BIN_DIR) mkdir $(BIN_DIR)
else
	@mkdir -p $(BIN_DIR)
endif
	go build -o $(BIN_DIR)/$(APP_NAME)



test: ## Run tests inside Docker container.
	docker run --rm -it $(APP_NAME) go test -v ./...

clean: ## Remove binaries
	@echo "Cleaning build artifacts..."
ifeq ($(OS),Windows_NT)
	@if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)

else
	@rm -rf $(BIN_DIR)
endif

docker: ## Build and run Docker containers using docker-compose
	docker-compose up --build
