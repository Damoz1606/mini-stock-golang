ifneq ("$(wildcard .env), "")
	include .env.makefile
	export
endif

CURDIR := $(shell powershell -Command "(Get-Location).Path")

SCHEMA_DIR := db/schema

MIGRATIONS_DIR := db/migrations

DB := $(DATABASE_TYPE)://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable
DB_TEMP := $(DATABASE_TYPE)://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_TMP)?sslmode=disable
TEST_DB := $(DATABASE_TYPE)://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(TEST_DATABASE)?sslmode=disable

.PHONY: run migrate lint test-unit test-coverage db-diff/% db-empty/% db-deploy db-hash db-status sqlc


# ----------------------
# Helpers
# ----------------------
define DATABASE_EXEC 
	$(CONTAINER_TOOL) exec -i $(DATABASE_CONTAINER) $(1)
endef

define RUN_ATLAS
	$(CONTAINER_TOOL) run --rm --network $(CONTAINER_NETWORK) -v $(CURDIR):/src -w /src $(IMAGE_ATLAS) $(1)
endef

define RUN_SQLC
	$(CONTAINER_TOOL) run --rm -v $(CURDIR):/src -w /src $(IMAGE_SQLC) generate
endef


# ----------------------
# Basic commands
# ----------------------
run:
	go run ./cmd/app

lint: 	
	golangci-lint run

mockery:
	$(CONTAINER_TOOL) run --rm -v $(CURDIR):/src -w /src vektra/mockery:3

test-unit:
	set CGO_ENABLED=1 && go test ./internal/... -v

test-coverage:
	set CGO_ENABLED=1 && go test -coverprofile=coverage.out -tags=testcoverage ./internal/...
	go tool cover -html=coverage.out -o coverage.html

test-setup-db:
	@echo "==> Dropping test DB if exists"
	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "DROP DATABASE IF EXISTS $(TEST_DATABASE);")

	@echo "==> Creating test DB"
	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "CREATE DATABASE $(TEST_DATABASE);")

	@echo "==> Applying migrations to $(TEST_DATABASE)"
	$(call RUN_ATLAS,migrate apply --dir "file://$(MIGRATIONS_DIR)" --url "$(TEST_DB)")

# ----------------------
# Database commands
# ----------------------
db-drop-temp:
	@echo "==> Dropping temp DB if exists"
	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "DROP DATABASE IF EXISTS $(DATABASE_TMP);")

db-create-temp:
	@echo "==> Creating temp DB"
	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "CREATE DATABASE $(DATABASE_TMP);")

db-diff/%: db-drop-temp db-create-temp
	@echo "==> Generating migration diff $*"
	$(call RUN_ATLAS,migrate diff $* --dir "file://$(MIGRATIONS_DIR)" --to "file://$(SCHEMA_DIR)/schema.pg.hcl" --dev-url "$(DB_TEMP)")
	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "DROP DATABASE IF EXISTS $(DATABASE_TMP);")

db-empty/%:
	$(call RUN_ATLAS,migrate new $* --dir "file://$(MIGRATIONS_DIR)")

db-deploy:
	$(call RUN_ATLAS,migrate apply --dir "file://$(MIGRATIONS_DIR)" --url "$(DB)")

db-hash:
	$(call RUN_ATLAS,migrate hash --dir "file://$(MIGRATIONS_DIR)")
			
db-status:
	$(call RUN_ATLAS,migrate status --dir "file://$(MIGRATIONS_DIR)" --url "$(DB)")

# ----------------------
# SQLC code generation
# ----------------------
sqlc: db-drop-temp db-create-temp
	@echo "==> Applying migrations to $(DATABASE_TMP)"
	$(call RUN_ATLAS,migrate apply --dir "file://$(MIGRATIONS_DIR)" --url "$(DB_TEMP)")

	@echo "==> Dumping schema"
	$(call DATABASE_EXEC,sh -c "pg_dump --schema-only --no-owner --no-privileges --no-reconnect --no-comments -U $(DATABASE_USER) -d $(DATABASE_TMP) | grep -v '^\\'") > $(SCHEMA_DIR)/schema.sql

	$(call DATABASE_EXEC,psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -c "DROP DATABASE IF EXISTS $(DATABASE_TMP);")

	@echo "==> Running sqlc"
	$(call RUN_SQLC)

sqlc-generate:
	$(call RUN_SQLC)