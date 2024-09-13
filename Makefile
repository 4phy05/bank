postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup1:
	migrate -path db/migration  -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrateup:
	migrate -path db/migration  -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown1:
	migrate -path db/migration  -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

migratedown:
	migrate -path db/migration  -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v D:\GoProject\Golang+Postgres+Docker:/src -w /src sqlc/sqlc:1.26.0 generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go SimpleBank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup1 migratedown1 migrateup migratedown sqlc test server mock