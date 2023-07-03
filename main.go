package main

import (
	"database/sql"
	"github.com/ismail118/simple-bank/api"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config error:", err)
	}

	conn, err := sql.Open(conf.DbDriver, conf.DbSource)
	if err != nil {
		log.Fatal("cannot connect db error:", err)
	}

	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot make token maker error:", err)
	}

	repo := repository.NewPostgresRepo(conn)
	store := repository.NewStore(conn)

	srv := api.NewServer(store, repo, tokenMaker, &conf)

	err = srv.Start(conf.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server error:", err)
	}
}
