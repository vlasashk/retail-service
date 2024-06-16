build-all:
	cd cart && make build
	cd loms && make build

.PHONY: run-all
run-all: export POSTGRES_PASSWORD=
run-all:
	docker compose up --force-recreate --build -d

down:
	docker compose down

up-cart-env:
	docker compose up -d --wait loms

up-loms-env:
	docker compose up -d --wait loms_db

.PHONY: cart-integration-test
cart-integration-test: up-cart-env
	cd cart && make integration-test
	docker compose down

.PHONY: loms-integration-test
loms-integration-test: up-loms-env
	cd loms && make integration-test
