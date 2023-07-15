package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"github.com/ismail118/simple-bank/worker"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
	"time"
)

var serverTest Server
var grpcServerTest *GrpcServer

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	repoMock := repository.NewPostgresRepoMock(nil)
	storeMock := repository.NewStoreMock(nil)
	taskDistributorMock := worker.NewRedisTaskDistributorMock(asynq.RedisClientOpt{})

	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Msgf("failed setup NewPasetoMaker error: %s", err)
	}

	serverTest = NewServer(storeMock, repoMock, tokenMaker, &config)

	// Grpc
	grpcServerTest = NewGrpcServer(storeMock, repoMock, tokenMaker, &config, taskDistributorMock)
	os.Exit(m.Run())
}
