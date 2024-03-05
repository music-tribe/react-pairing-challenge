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

test: mocks start
	cd api/ && go test ./... -v
	make stop

mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	cd api && go generate ./...

docs:
	cd api && swag init -g *.go --output docs/features-api
