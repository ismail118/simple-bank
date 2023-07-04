package main

import (
	"database/sql"
	"github.com/ismail118/simple-bank/api"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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

	// run http server
	//runGinServer(store, repo, tokenMaker, conf)

	// run grpc server
	runGrpcServer(store, repo, tokenMaker, conf)

}

func runGrpcServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config) {
	server := api.NewGrpcServer(store, repo, tokenMaker, &conf)

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	// reflection.Register(grpcServer) is optional but recommended
	//to let know client what services RPCs available on the server and how to call them (documentation)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", conf.GrpcServerAddr)
	if err != nil {
		log.Fatalf("error listen grpc addr:%s", conf.GrpcServerAddr)
	}

	log.Println("start gRPC server at ", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start grp server error:%s", err)
	}
}

func runGinServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config) {
	srv := api.NewServer(store, repo, tokenMaker, &conf)

	err := srv.Start(conf.HttpServerAddr)
	if err != nil {
		log.Fatal("cannot start server error:", err)
	}
}
