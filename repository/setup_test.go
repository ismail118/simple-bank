package repository

import (
	"context"
	"github.com/ismail118/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err)
	}
	dbpool, err := pgxpool.New(context.Background(), conf.DbSource)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	repo := NewPostgresRepo(dbpool)
	testRepo = repo

	store := NewStore(dbpool)
	testStore = store

	os.Exit(m.Run())
}
