package main

import (
	"database/sql"
	"github.com/ismail118/simple-bank/api"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(conf.DbDriver, conf.DbSource)
	if err != nil {
		log.Fatal("cannot connect db error:", err)
	}

	repo := repository.NewPostgresRepo(conn)
	store := repository.NewStore(conn)

	srv := api.NewServer(store, repo)

	err = srv.Start(conf.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server error:", err)
	}
}
