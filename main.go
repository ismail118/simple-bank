package main

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/ismail118/simple-bank/api"
	"github.com/ismail118/simple-bank/mail"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"github.com/ismail118/simple-bank/worker"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config error:%s", err)
	}

	// run db migration
	runDBMigration(conf.MigrationURL, conf.DbSource)

	conn, err := sql.Open(conf.DbDriver, conf.DbSource)
	if err != nil {
		log.Fatal().Msgf("cannot connect db error:%s", err)
	}

	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Msgf("cannot make token maker error:%s", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: conf.RedisAddr,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	mailer := mail.NewGmailSender(conf.EmailSenderName, conf.EmailSenderAddress, conf.EmailSenderPassword)

	repo := repository.NewPostgresRepo(conn)
	store := repository.NewStore(conn)

	// run task processor
	go runTaskProcessor(redisOpt, store, mailer, conf.GatewayServerAddr)

	// run grpc server
	go runGrpcServer(store, repo, tokenMaker, conf, taskDistributor)

	go runGatewayServer(store, repo, tokenMaker, conf, taskDistributor)

	// run http server
	runGinServer(store, repo, tokenMaker, conf)

}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msgf("cannot create db migration err:%s", err)
	}

	err = migration.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal().Msgf("cannot run migration up err:%s", err)
		}
	}

	log.Info().Msg("success run db migration")
}

func runGrpcServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config, taskDistributor worker.TaskDistributor) {
	server := api.NewGrpcServer(store, repo, tokenMaker, &conf, taskDistributor)

	// middleware/interceptor logger
	grpcInterceptor := grpc.UnaryInterceptor(api.GrpcInterceptorLogger)

	grpcServer := grpc.NewServer(grpcInterceptor)
	pb.RegisterSimpleBankServer(grpcServer, server)

	// reflection.Register(grpcServer) is optional but recommended
	//to let know client what services RPCs available on the server and how to call them (documentation)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", conf.GrpcServerAddr)
	if err != nil {
		log.Fatal().Msgf("error listen grpc addr:%s", conf.GrpcServerAddr)
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("cannot start grp server error:%s", err)
	}
}

func runGatewayServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config, taskDistributor worker.TaskDistributor) {
	server := api.NewGrpcServer(store, repo, tokenMaker, &conf, taskDistributor)

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
		log.Fatal().Msgf("cannot register handler server err:%s", err)
	}

	// reroute http mux to grpcMux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", conf.GatewayServerAddr)
	if err != nil {
		log.Fatal().Msgf("error listen gateway addr:%s", conf.GrpcServerAddr)
	}

	log.Info().Msgf("start http gateway server at %s", listener.Addr().String())
	handler := api.HttpGatewayInterceptorLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msgf("cannot start http gateway server error:%s", err)
	}
}

func runGinServer(store repository.Store, repo repository.Repository, tokenMaker token.Maker, conf util.Config) {
	srv := api.NewServer(store, repo, tokenMaker, &conf)

	err := srv.Start(conf.HttpServerAddr)
	if err != nil {
		log.Fatal().Msgf("cannot start server error:%s", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store repository.Store, mailer mail.SenderEmail, gatewaySeverAddress string) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer, gatewaySeverAddress)
	log.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
