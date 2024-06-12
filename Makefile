build-all:
	cd cart && make build
	cd loms && make build


run-all:
	docker-compose up --force-recreate --build -d

down:
	docker-compose down

