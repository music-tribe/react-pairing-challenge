start:
	docker compose -f api/docker-compose.yaml up

start_no_logs:
	docker compose -f api/docker-compose.yaml up -d --remove-orphans

stop: 
	docker compose -f api/docker-compose.yaml down