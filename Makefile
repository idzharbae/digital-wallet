build:
	@echo " >> building binaries"
	@go build -v -o bin/digital-wallet src/cmd/main.go

run: build
	@./bin/digital-wallet

migrate:
	@migrate -path database/migration/ -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate-down:
	@migrate -path database/migration/ -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

generate-mock:
	@go generate ./...

test:
	@go test ./...
