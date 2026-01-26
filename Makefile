tidy:
	@go mod tidy

up-local:
	@docker compose -f docker/docker-compose.local.yaml up -d

sqlc:
	@docker compose -f docker/docker-compose.local.yaml run --rm local_sqlc generate
