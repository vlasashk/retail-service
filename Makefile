build-all:
	cd cart && make build
	cd loms && make build

run-all:
	docker compose up --force-recreate --build -d

down:
	docker compose down

up-cart-env:
	docker compose up -d --wait loms


.PHONY: cart-integration-test
cart-integration-test: up-cart-env
	cd cart && make integration-test
	docker compose down

.PHONY: loms-integration-test
loms-integration-test:
	cd loms && make integration-test
	docker compose down
