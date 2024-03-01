start_with_logs:
	docker compose -f api/docker-compose.yaml up

start:
	docker compose -f api/docker-compose.yaml up -d --remove-orphans

stop: 
	docker compose -f api/docker-compose.yaml down

build:
	docker compose -f api/docker-compose.yaml build