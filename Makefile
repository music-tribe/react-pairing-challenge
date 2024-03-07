start_with_logs:
	docker compose -f api/docker-compose.yml up

start:
	docker compose -f api/docker-compose.yml up -d --remove-orphans

start_db:
	docker compose -f api/docker-compose.yml run database

stop: 
	docker compose -f api/docker-compose.yml down

build: mocks docs
	docker compose -f api/docker-compose.yml build

test: mocks start
	cd api/ && go test ./... -v
	make stop

mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	cd api && go generate ./...

docs:
	cd api && go install github.com/swaggo/swag/cmd/swag@latest && swag init -g *.go --output docs/features-api