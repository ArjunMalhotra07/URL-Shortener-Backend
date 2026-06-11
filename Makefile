tidy:
	@go mod tidy

up-local:
	@docker compose -f docker/docker-compose.local.yaml up -d

sqlc:
	@docker compose -f docker/docker-compose.local.yaml run --rm local_sqlc generate


goose_up:
	@goose -dir db/migrations postgres "$DB_DSN" up


totp:
	@oathtool --totp -b NE67GTLV6C32PAW24YFKX2E5X4KI6WNZ