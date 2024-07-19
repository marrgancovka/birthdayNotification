include .env

migrate-lib:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: migrate-lib
	migrate -path schema -database "postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path schema -database "postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

run:
	docker-compose build
	docker-compose up -d

stop:
	docker-compose down
