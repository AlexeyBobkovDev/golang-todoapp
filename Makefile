include .env
export 

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Do you want to clean up all Postgres volumes? DANGEROUS! [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf $(PROJECT_ROOT)/out/pgdata && \
		echo "All files were removed"; \
	else \
		echo "Clean up was discarded"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder
	

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Missing required parameter (name). Example of using this target: make migrate-create name=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate\
		create \
		-ext sql \
		-dir /migrations \
		-seq \
		"$(name)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter (action). Example: make migrate-action action=up" \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

todoapp-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run $(PROJECT_ROOT)/cmd/todoapp/main.go