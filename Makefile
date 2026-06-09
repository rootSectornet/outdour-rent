.PHONY: dev dev-api dev-web migrate migrate-down seed

run-dev:
	cd apps/api && go run cmd/server/main.go

dev:
	docker compose -f deployments/docker-compose.yml -f deployments/docker-compose.dev.yml up

dev-api:
	cd apps/api && air

dev-web:
	cd apps/web && npm run dev

migrate:
	cd apps/api && go run cmd/migrate/main.go up

migrate-down:
	cd apps/api && go run cmd/migrate/main.go down

seed:
	cd apps/api && go run cmd/migrate/main.go seed
