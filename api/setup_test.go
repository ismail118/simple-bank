package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/repository"
	"os"
	"testing"
)

var serverTest Server

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	repoMock := repository.NewPostgresRepoMock(nil)
	storeMock := repository.NewStoreMock(nil)

	serverTest = NewServer(storeMock, repoMock)
	os.Exit(m.Run())
}
