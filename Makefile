generate:
	sqlc generate
dev:
	fresh
migrate:
	cd sql/schema; goose postgres postgres://postgres:pass@localhost:5432/data up
dev-db:
	docker compose up db -d
prod:
	docker compose up --build
tidy:
	go mod tidy