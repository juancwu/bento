DB_CONN ?= $(shell cat .env | grep BENTO_DB_CONN | sed 's/^[^=]*=//')
MIGRATION_FOLDER ?= ./migrations

migrate-up:
	@libsql-migrate up --url "$(DB_CONN)" --path "$(MIGRATION_FOLDER)"

migrate-down:
	@libsql-migrate down --url "$(DB_CONN)" --path "$(MIGRATION_FOLDER)"

migrate-new:
	@libsql-migrate gen $(NAME) --path $(MIGRATION_FOLDER)
