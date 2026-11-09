.PHONY: up down test build
up:
	docker compose up --build
down:
	docker compose down
test:
	cd backend && go test ./...
build:
	cd backend && go build ./cmd/api
	cd frontend && npm run build
