MIGRATION_FOLDER ?= ./migrations

run:
	@go run ./cmd/bento

dev:
	@air

migrate-up:
	@go run ./migrate.go up

migrate-down:
	@go run ./migrate.go down

migrate-new:
	@migrate create -ext sql -dir $(MIGRATION_FOLDER) -seq $(NAME)
