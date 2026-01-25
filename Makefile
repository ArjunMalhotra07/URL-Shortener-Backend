run:
	@go run ./cmd/server

build:
	@go build -o bin/server ./cmd/server

docker-build:
	@docker build -t url-shortner-backend -f docker/Dockerfile .

docker-up:
	@docker compose -f docker/docker-compose.yaml up --build

tidy:
	@go mod tidy
