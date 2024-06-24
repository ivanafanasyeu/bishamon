docker-dev-up:
	@cp .env.dev ./docker/dev/.env & docker-compose -f ./docker/dev/docker-compose.dev.yaml up -d

docker-dev-stop:
	@docker-compose -f ./docker/dev/docker-compose.dev.yaml stop

docker-dev-down:
	@docker-compose -f ./docker/dev/docker-compose.dev.yaml down

dev:
	@air

build:
	@templ generate
	@go build -o bin/bishamon cmd/server/main.go

run: build
	@./bin/app

test:
	@go test -v ./...