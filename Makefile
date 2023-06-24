install_golang_migration:
	echo "install golang migration"
	brew install golang-migrate

create_migration:
	echo "create migration file"
	migrate create -ext sql -dir db/migration -seq init_schema

create_db:
	echo "create database"
	docker exec -it postgres-db-1 createdb --username=postgres --owner=postgres simple_bank

drop_db:
	echo "drop database"
	docker exec -it postgres-db-1 dropdb --username=postgres simple_bank

migrate_up:
	echo "up migration"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	echo "down migration"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate_up1:
	echo "up migration 1"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate_down1:
	echo "down migration 1"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

test:
	go test -v -cover ./...

