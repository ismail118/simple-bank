package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var tokenMakerTest token.Maker

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// taskDistributorMock := worker.NewRedisTaskDistributorMock(asynq.RedisClientOpt{})

	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Msgf("failed setup NewPasetoMaker error: %s", err)
	}

	tokenMakerTest = tokenMaker

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, repo repository.Repository, store repository.Store, tokenMaker token.Maker) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server := NewServer(store, repo, tokenMaker, &config)

	return &server
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, _, err := tokenMaker.CreateToken(username, duration)
	assert.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account models.Account) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var gotAccount models.Account
	err = json.Unmarshal(data, &gotAccount)
	assert.NoError(t, err)

	assert.Equal(t, account, gotAccount)
}

type eqCreateUserTxParamsMatcher struct {
	arg      models.Users
	password string
	expected     models.Users
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(models.Users)
	if !ok {
		return false
	}

	err := util.ComparePassword(actualArg.HashedPassword, expected.password)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg, actualArg) {
		return false
	}

	return true
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg models.Users, password string, expected models.Users) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, expected}
}