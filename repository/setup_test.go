package repository

import (
	"database/sql"
	"github.com/ismail118/simple-bank/util"
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
var testStore Store

func TestMain(m *testing.M) {
	conf, err := util.LoadConfig("../.")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := sql.Open(conf.DbDriver, conf.DbSource)
	if err != nil {
		log.Fatal(err)
	}

	repo := NewPostgresRepo(conn)
	testRepo = repo

	store := NewStore(conn)
	testStore = store

	os.Exit(m.Run())
}
