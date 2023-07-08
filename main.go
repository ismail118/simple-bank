package main

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ismail118/simple-bank/api"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
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

	// run grpc server
	go runGrpcServer(store, repo, tokenMaker, conf)

	go runGatewayServer(store, repo, tokenMaker, conf)

	// run http server
	runGinServer(store, repo, tokenMaker, conf)

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

func runGatewayServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config) {
	server := api.NewGrpcServer(store, repo, tokenMaker, &conf)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("cannot register handler server err:%s", err)
	}

	// reroute http mux to grpcMux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", conf.GatewayServerAddr)
	if err != nil {
		log.Fatalf("error listen gateway addr:%s", conf.GrpcServerAddr)
	}

	log.Println("start http gateway server at ", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start http gateway server error:%s", err)
	}
}

func runGinServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config) {
	srv := api.NewServer(store, repo, tokenMaker, &conf)

	err := srv.Start(conf.HttpServerAddr)
	if err != nil {
		log.Fatal("cannot start server error:", err)
	}
}
