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
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_up_docker:
	echo "up migration"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	docker exec -it simple-bank-simplebank-1 /app/migrate -path db/migration -database "postgres://postgres:postgres@postgres-service:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	echo "down migration"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate_up1:
	echo "up migration 1"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate_down1:
	echo "down migration 1"
	# -database "driver://user:password@host:port/db_name?sslmode=disable"
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

test:
	go test -v -cover ./...

build_docker_image:
	echo "build docker image from Dockerfile"
	docker build -t ismail118/simplebank:1.0.0 -f Dockerfile .

run_docker_container:
	echo "run docker simplebank"
	docker run --rm --name simplebank -p 8080:8080 -e GIN_MODE=release ismail118/simplebank:1.0.0

compose_up_build:
	docker-compose up --build