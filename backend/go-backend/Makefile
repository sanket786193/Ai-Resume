.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset run build deps

# Database migration commands
migrate-up:
	@go run cmd/migrate/main.go -command=up

migrate-down:
	@go run cmd/migrate/main.go -command=down

migrate-status:
	@go run cmd/migrate/main.go -command=status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@go run cmd/migrate/main.go -command=create $(name)

migrate-reset:
	@go run cmd/migrate/main.go -command=down
	@go run cmd/migrate/main.go -command=up

# Application commands
run:
	@go run main.go

build:
	@go build -o bin/resume-backend main.go

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Docker commands
docker-build:
	@docker-compose build

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-restart:
	@docker-compose restart

docker-clean:
	@docker-compose down -v
	@docker-compose rm -f

docker-migrate:
	@docker-compose run --rm migrate -command=up

docker-migrate-status:
	@docker-compose run --rm migrate -command=status
