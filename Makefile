init_db_container:
	docker run --name bankingSampleDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1 -d postgres

start_container:
	docker start bankingSampleDB

stop_container:
	docker stop bankingSampleDB

createdb_simple_bank:
	docker exec -it bankingSampleDB createdb --username=root --owner=root simple_bank

createdb_test_simple_bank:
	docker exec -it bankingSampleDB createdb --username=root --owner=root test_simple_bank

dropdb_simple_bank:
	docker exec -it bankingSampleDB dropdb simple_bank

dropdb_test_simple_bank:
	docker exec -it bankingSampleDB dropdb test_simple_bank

migrate_up:
	migrate -path db/migration -database "postgresql://root:1@localhost:5432/simple_bank?sslmode=disable" -verbose up
	migrate -path db/migration -database "postgresql://root:1@localhost:5432/test_simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:1@localhost:5432/simple_bank?sslmode=disable" -verbose down
	migrate -path db/migration -database "postgresql://root:1@localhost:5432/test_simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: init_db_container start_container stop_container
.PHONY: createdb_simple_bank createdb_test_simple_bank dropdb_simple_bank dropdb_test_simple_bank
.PHONY: migrate_up migrate_down
.PHONY: sqlc test