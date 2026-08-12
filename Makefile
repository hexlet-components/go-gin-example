MIGRATOR=go run main.go migrate

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Download dependencies
	go mod download

tidy: ## Tidy up dependencies, format code, and run vet
	go mod tidy
	go fmt ./...
	go vet ./...

lint: ## Check formatting and run vet without changing files
	@files=$$(gofmt -l .); if [ -n "$$files" ]; then echo "gofmt is needed for:"; echo "$$files"; exit 1; fi
	go vet ./...

dev: ## Run the API server in development mode
	go run main.go api

.PHONY: install lint test
test: ## Run all tests
	go test -v ./...

db-migrate: ## Run database migrations
	$(MIGRATOR) up

db-rollback: ## Rollback database migrations
	$(MIGRATOR) down

db-reset: ## Reset database to initial state
	$(MIGRATOR) reset

db-status: ## Show database migration status
	$(MIGRATOR) status

db-generate: ## Generate database code using sqlc
	go tool sqlc generate