export PG_PASS=postgres
export PG_USER=postgres
export PG_DB=postgres
export PG_HOST=pg_master:5432

.PHONY: build-all
build-all:
	cd cart && make build
	cd loms && make build

.PHONY: run-all
run-all:
	docker compose up --force-recreate --build -d

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: up-cart-env
up-cart-env:
	docker compose up -d --wait loms

.PHONY: up-loms-env
up-loms-env:
	docker compose up -d --wait pg_master pg_slave migration

.PHONY: cart-integration-test
cart-integration-test: export STOCKS_DB_PASSWORD=$(PG_PASS)
cart-integration-test: export ORDERS_DB_PASSWORD=$(PG_PASS)
cart-integration-test: up-cart-env
	cd cart && make integration-test
	docker compose down

.PHONY: loms-integration-test
loms-integration-test: export STOCKS_DB_PASSWORD=$(PG_PASS)
loms-integration-test: export ORDERS_DB_PASSWORD=$(PG_PASS)
loms-integration-test: up-loms-env
	cd loms && make integration-test
	docker compose down
