package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ismail118/simple-bank/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_getAccount(t *testing.T) {
	testCases := []struct {
		Name               string
		AccountID          int64
		setupRequest       func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectedStatusCode int
	}{
		{
			Name:      "Accepted",
			AccountID: 2,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:      "BadReq",
			AccountID: 0,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:      "ServerError",
			AccountID: 1001,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:      "NotFound",
			AccountID: 1,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectedStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectedStatusCode, rr.Code)
		}
	}
}

func Test_createAccount(t *testing.T) {
	testCase := []struct {
		Name                  string
		ReqBody               map[string]interface{}
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"currency": "WRONG",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", reqBody)
		req.Header.Set("Content-Type", "application/json")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_listAccount(t *testing.T) {
	testCase := []struct {
		Name                  string
		QueryParams           string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:        "Accepted",
			QueryParams: "page=1&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:        "ServerError",
			QueryParams: "page=1001&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_updateAccount(t *testing.T) {
	testCase := []struct {
		Name                  string
		ReqBody               map[string]interface{}
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"id":       2,
				"balance":  100,
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"balance":  100,
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name: "ServerError",
			ReqBody: map[string]interface{}{
				"id":       1001,
				"balance":  100,
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "ServerErrorUpdate",
			ReqBody: map[string]interface{}{
				"id":       3,
				"balance":  100,
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"id":       1000,
				"balance":  100,
				"currency": "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPut, "/accounts", reqBody)
		req.Header.Set("Content-Type", "application/json")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_deleteAccount(t *testing.T) {
	testCases := []struct {
		Name               string
		AccountID          int64
		setupRequest       func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectedStatusCode int
	}{
		{
			Name:      "Accepted",
			AccountID: 2,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:      "BadReq",
			AccountID: 0,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:      "ServerError",
			AccountID: 1001,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:      "ServerErrorDelete",
			AccountID: 3,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:      "NotFound",
			AccountID: 1,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/accounts/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)
		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectedStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectedStatusCode, rr.Code)
		}
	}
}

func Test_getEntry(t *testing.T) {
	testCases := []struct {
		Name               string
		AccountID          int64
		setupRequest       func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectedStatusCode int
	}{
		{
			Name:      "Accepted",
			AccountID: 2,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:      "BadReq",
			AccountID: 0,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:      "ServerError",
			AccountID: 1001,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:      "NotFound",
			AccountID: 1,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entries/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)
		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectedStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectedStatusCode, rr.Code)
		}
	}
}

func Test_listEntries(t *testing.T) {
	testCase := []struct {
		Name                  string
		QueryParams           string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:        "Accepted",
			QueryParams: "account_id=2&page=1&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:        "ServerError",
			QueryParams: "account_id=2&page=1001&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entries?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_getTransfer(t *testing.T) {
	testCases := []struct {
		Name               string
		TransferID         int64
		setupRequest       func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectedStatusCode int
	}{
		{
			Name:       "Accepted",
			TransferID: 2,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:       "BadReq",
			TransferID: 0,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:       "ServerError",
			TransferID: 1001,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:       "NotFound",
			TransferID: 1,
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfer/%d", tc.TransferID), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectedStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectedStatusCode, rr.Code)
		}
	}
}

func Test_listTransfer(t *testing.T) {
	testCase := []struct {
		Name                  string
		QueryParams           string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:        "Accepted",
			QueryParams: "from_account_id=1&to_account_id=2&page=1&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:        "ServerError",
			QueryParams: "from_account_id=1&to_account_id=2&page=1001&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfer?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_transfer(t *testing.T) {
	testCase := []struct {
		Name                  string
		ReqBody               map[string]interface{}
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"from_account_id": 2,
				"to_account_id":   3,
				"amount":          10,
				"currency":        "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"from_account_id": 0,
				"to_account_id":   0,
				"amount":          0,
				"currency":        "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name: "ServerError",
			ReqBody: map[string]interface{}{
				"from_account_id": 1001,
				"to_account_id":   3,
				"amount":          10,
				"currency":        "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "ServerError",
			ReqBody: map[string]interface{}{
				"from_account_id": 2,
				"to_account_id":   2001,
				"amount":          10,
				"currency":        "USD",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPost, "/transfer", reqBody)
		req.Header.Set("Content-Type", "application/json")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_createUsers(t *testing.T) {
	testCases := []struct {
		Name                  string
		ReqBody               map[string]interface{}
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"username":  "user",
				"password":  "secret",
				"full_name": "david jones",
				"email":     "notexists@gmail.com",
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
	}

	for _, tc := range testCases {
		reqBody, _ := json.Marshal(tc.ReqBody)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_getUsers(t *testing.T) {
	testCases := []struct {
		Name                  string
		Username              string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:     "Accepted",
			Username: "some-user",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name:     "NotFound",
			Username: "user",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tc.Username), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_listUsers(t *testing.T) {
	testCase := []struct {
		Name                  string
		QueryParams           string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:        "Accepted",
			QueryParams: "page=1&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:        "ServerError",
			QueryParams: "page=1001&size=10",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_updateUsers(t *testing.T) {
	testCase := []struct {
		Name                  string
		ReqBody               map[string]interface{}
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"username":  "some-user",
				"full_name": "ismail",
				"email":     "notexists@gmail.com",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"username": "ismail",
				"email":    "some@email.com",
			},
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPut, "/users", reqBody)
		req.Header.Set("Content-Type", "application/json")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}

func Test_deleteUsers(t *testing.T) {
	testCases := []struct {
		Name                  string
		Username              string
		setupRequest          func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		ExpectationStatusCode int
	}{
		{
			Name:     "Accepted",
			Username: "some-user",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name:     "NotFound",
			Username: "user",
			setupRequest: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken("some-user", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}
				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			ExpectationStatusCode: http.StatusNotFound,
		},
	}

	tokenMaker := serverTest.tokenMaker

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", tc.Username), nil)
		req.Header.Set("Content-Type", "text/plain")

		tc.setupRequest(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}
