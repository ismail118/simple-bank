package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"
)

var testRepo Repository
var testStore *Store

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	repo := NewPostgresRepo(conn)
	testRepo = repo

	store := NewStore(conn)
	testStore = store

	os.Exit(m.Run())
}
