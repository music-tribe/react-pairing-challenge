start_with_logs:
	docker compose -f api/docker-compose.yaml up

start:
	docker compose -f api/docker-compose.yaml up -d --remove-orphans

start_db:
	docker compose -f api/docker-compose.yaml run database

stop: 
	docker compose -f api/docker-compose.yaml down

build:
	docker compose -f api/docker-compose.yaml build

test: start
	cd api/ && go test ./... -v
	make stop
