postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ultramegasecret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root neobank

dropdb:
	docker exec -it postgres12 dropdb neobank

migrateup:
	migrate -path db/migrations -database "postgresql://root:ultramegasecret@localhost:5432/neobank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:ultramegasecret@localhost:5432/neobank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:     
	go test -v -cover ./...

serve:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
