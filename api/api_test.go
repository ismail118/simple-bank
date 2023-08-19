package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)

	account := util.RandomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repository.MockRepository)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(repo *repository.MockRepository) {
				repo.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := repository.NewMockStore(ctrl)

			repo := repository.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(t, repo, store, tokenMakerTest)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)
	account := util.RandomAccount(user.Username)
	account.Balance = 0

	testCases := []struct {
		name string
		postData gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repository.MockRepository)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	} {
		{
			name: "ok",
			postData: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(repo *repository.MockRepository) {
				acc := models.Account{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}

				repo.
				EXPECT().
				InsertAccount(gomock.Any(), gomock.Eq(acc)).
				Times(1).
				Return(account.ID, nil)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, rr.Code)
				requireBodyMatchAccount(t, rr.Body, account)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := repository.NewMockStore(ctrl)

			repo := repository.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(t, repo, store, tokenMakerTest)

			rr := httptest.NewRecorder()

			data, err := json.Marshal(tc.postData)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(data))

			assert.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(rr, request)
			tc.checkResponse(t, rr)
		})
	}
}

func TestListAccounts(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)
	accounts := []*models.Account{}

	testCases := []struct {
		name string
		params url.Values
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repository.MockRepository)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	} {
		{
			name: "ok",
			params: url.Values{
				"page" : {"1"},
				"size" : {"10"},
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMakerTest, authorizationBearer, user.Username, time.Minute)
			},
			buildStubs: func(repo *repository.MockRepository) {
				repo.
				EXPECT().
				GetListAccounts(
					gomock.Any(), 
					gomock.Eq(user.Username), 
					gomock.Eq(10), 
					gomock.Eq(0),
				).
				Times(1).
				Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, rr.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := repository.NewMockStore(ctrl)

			repo := repository.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(t, repo, store, tokenMakerTest)

			rr := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?%s", tc.params.Encode())
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, req, tokenMakerTest)

			server.router.ServeHTTP(rr, req)

			tc.checkResponse(t, rr)
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)
	account := util.RandomAccount(user.Username)

	testCases := []struct {
		name string
		postData gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repository.MockRepository, postData gin.H)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	} {
		{
			name: "ok",
			postData: gin.H{
				"id": account.ID,
				"balance": util.RandomBalance(),
				"currency": util.RandomCurrency(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(repo *repository.MockRepository, postData gin.H) {
				repo.
				EXPECT().
				GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)

				account.Balance = postData["balance"].(int64)
				account.Currency = postData["currency"].(string)

				repo.
				EXPECT().
				UpdateAccount(gomock.Any(), gomock.Eq(account)).
				Times(1).
				Return(nil)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, rr.Code)
				requireBodyMatchAccount(t, rr.Body, account)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := repository.NewMockStore(ctrl)

			repo := repository.NewMockRepository(ctrl)
			tc.buildStubs(repo, tc.postData)

			server := newTestServer(t, repo, store, tokenMakerTest)

			rr := httptest.NewRecorder()

			data, err := json.Marshal(tc.postData)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/accounts", bytes.NewReader(data))

			assert.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(rr, request)
			tc.checkResponse(t, rr)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)

	account := util.RandomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repository.MockRepository)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(repo *repository.MockRepository) {
				repo.EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)

				repo.EXPECT().
				DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := repository.NewMockStore(ctrl)

			repo := repository.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(t, repo, store, tokenMakerTest)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}