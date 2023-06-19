package main

import (
	"database/sql"
	"github.com/ismail118/simple-bank/api"
	"github.com/ismail118/simple-bank/repository"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect db error:", err)
	}

	repo := repository.NewPostgresRepo(conn)
	store := repository.NewStore(conn)

	srv := api.NewServer(store, repo)

	err = srv.Start(serverAddr)
	if err != nil {
		log.Fatal("cannot start server error:", err)
	}
}
