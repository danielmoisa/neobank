# Define the target to load environment variables from the .env file
include .env

# Define environment variables
DB_URL := postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Run PostgreSQL container
postgres:
	docker run --name neobank-postgres-1 -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:12-alpine

# Create the database
createdb:
	docker exec -it neobank-postgres-1 createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

# Drop the database
dropdb:
	docker exec -it neobank-postgres-1 dropdb $(DB_NAME)

# Migrate up
migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

# Migrate down
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

# Generate SQLC code
sqlc:
	sqlc generate

# Run tests
test:     
	go test -v -cover ./...

# Serve the application
serve:
	go run main.go

# Generate mock DB
mock:
	mockgen -package mockdb -destination db/mocks/store.go github.com/danielmoisa/neobank/db/sqlc Store

# Generate Swagger docs
docs:
	swag init

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test serve mock docs
