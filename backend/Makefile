.PHONY: run build test clean db-up db-down db-logs migrate-status

run:
	go run cmd/app/main.go

build:
	go build -o bin/url-shortener cmd/app/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf bin

db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres

db-shell:
	docker compose exec postgres psql -U $${DB_USER:-postgres} -d $${DB_NAME:-url_shortener}

migrate-status:
	@migrate -source file://internal/database/migrations \
		-database "postgres://$${DB_USER:-postgres}:$${DB_PASSWORD:-postgres}@localhost:$${DB_PORT:-5432}/$${DB_NAME:-url_shortener}?sslmode=disable" \
		status || true