include .env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

export PROJECT_ROOT=$(shell pwd)


# ==========================ENV==========================
.PHONY: env-up
env-up:
	@docker compose up -d todoapp-postgres

.PHONY: env-down
env-down:
	@docker compose down todoapp-postgres

.PHONY: env-cleanup
env-cleanup:
	@read -p "Do you want to clean up all Postgres volumes? DANGEROUS! [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf $(PROJECT_ROOT)/out/pgdata && \
		echo "All files were removed"; \
	else \
		echo "Clean up was discarded"; \
	fi

.PHONY: env-port-forward
env-port-forward:
	@docker compose up -d port-forwarder

.PHONY: env-port-close
env-port-close:
	@docker compose down port-forwarder


# ==========================MIGRATE==========================
.PHONY: migrate-create
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

.PHONY: migrate-up
migrate-up:
	@$(MAKE) migrate-action action=up

.PHONY: migrate-down
migrate-down:
	@$(MAKE) migrate-action action=down

.PHONY: migrate-action
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter (action). Example: make migrate-action action=up" \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"


# ==========================APP==========================
.PHONY: todoapp-run
todoapp-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	export POSTGRES_HOST=localhost && \
	go run $(PROJECT_ROOT)/cmd/todoapp/main.go

.PHONY: todoapp-deploy
todoapp-deploy: env-up   env-port-forward   migrate-up   swagger-gen   env-port-close
	@docker compose up -d --build todoapp

.PHONY: todoapp-undeploy
todoapp-undeploy:
	@docker compose down todoapp


# ==========================SWAGGER==========================
.PHONY: swagger-gen
CMD = docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal
swagger-gen: swagger-fmt
	@$(CMD) || $(CMD) --parseDependency


.PHONY: swagger-fmt
swagger-fmt:
	@$(MAKE) swagger-action action=fmt

.PHONY: swagger-action
swagger-action:
	@if [ -z "$(action)" ]; then  \
		echo "Missed required parameter `action`" && \
		exit 1; \
	fi; \
	docker compose run --rm swagger $(action) \
		--dir ./cmd,./internal


# ==========================DOCKER==========================
.PHONY: ps
ps:
	@docker compose ps


# ==========================LOGS==========================
.PHONY: logs-cleanup
logs-cleanup:
	@read -p "Do you want to clean up all logs? DANGEROUS! [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf $(PROJECT_ROOT)/out/logs && \
		echo "All files were removed"; \
	else \
		echo "Clean up was discarded"; \
	fi


# ==========================GOLANG==========================

.PHONY: deps
deps:
	@go mod download

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: fmt
fmt:
	@go fmt ./...
