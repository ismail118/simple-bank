package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getAccount(t *testing.T) {
	testCases := []struct {
		Name               string
		AccountID          int64
		ExpectedStatusCode int
	}{
		{
			Name:               "Accepted",
			AccountID:          2,
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:               "BadReq",
			AccountID:          0,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "ServerError",
			AccountID:          1001,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "NotFound",
			AccountID:          1,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"owner":    "ismail",
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"owner":    "ismail",
				"currency": "WRONG",
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name: "ServerError",
			ReqBody: map[string]interface{}{
				"owner":    "test-error-db",
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", reqBody)
		req.Header.Set("Content-Type", "application/json")

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
		ExpectationStatusCode int
	}{
		{
			Name:                  "Accepted",
			QueryParams:           "page=1&size=10",
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name:                  "BadRequest",
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:                  "ServerError",
			QueryParams:           "page=1001&size=10",
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		ExpectationStatusCode int
	}{
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"id":       2,
				"owner":    "ismail",
				"balance":  100,
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			ReqBody: map[string]interface{}{
				"owner":    "ismail",
				"balance":  100,
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name: "ServerError",
			ReqBody: map[string]interface{}{
				"id":       1001,
				"owner":    "ismail",
				"balance":  100,
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "ServerErrorUpdate",
			ReqBody: map[string]interface{}{
				"id":       3,
				"owner":    "ismail",
				"balance":  100,
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusInternalServerError,
		},
		{
			Name: "Accepted",
			ReqBody: map[string]interface{}{
				"id":       1000,
				"owner":    "ismail",
				"balance":  100,
				"currency": "USD",
			},
			ExpectationStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPut, "/accounts", reqBody)
		req.Header.Set("Content-Type", "application/json")

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
		ExpectedStatusCode int
	}{
		{
			Name:               "Accepted",
			AccountID:          2,
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:               "BadReq",
			AccountID:          0,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "ServerError",
			AccountID:          1001,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "ServerErrorDelete",
			AccountID:          3,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "NotFound",
			AccountID:          1,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/accounts/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		ExpectedStatusCode int
	}{
		{
			Name:               "Accepted",
			AccountID:          2,
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:               "BadReq",
			AccountID:          0,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "ServerError",
			AccountID:          1001,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "NotFound",
			AccountID:          1,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entries/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		ExpectationStatusCode int
	}{
		{
			Name:                  "Accepted",
			QueryParams:           "account_id=1&page=1&size=10",
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name:                  "BadRequest",
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:                  "ServerError",
			QueryParams:           "account_id=1&page=1001&size=10",
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/entries?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		AccountID          int64
		ExpectedStatusCode int
	}{
		{
			Name:               "Accepted",
			AccountID:          2,
			ExpectedStatusCode: http.StatusAccepted,
		},
		{
			Name:               "BadReq",
			AccountID:          0,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "ServerError",
			AccountID:          1001,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "NotFound",
			AccountID:          1,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfer/%d", tc.AccountID), nil)
		req.Header.Set("Content-Type", "text/plain")

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
		ExpectationStatusCode int
	}{
		{
			Name:                  "Accepted",
			QueryParams:           "from_account_id=1&to_account_id=2&page=1&size=10",
			ExpectationStatusCode: http.StatusAccepted,
		},
		{
			Name:                  "BadRequest",
			ExpectationStatusCode: http.StatusBadRequest,
		},
		{
			Name:                  "ServerError",
			QueryParams:           "from_account_id=1&to_account_id=2&page=1001&size=10",
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfer?%s", tc.QueryParams), nil)
		req.Header.Set("Content-Type", "text/plain")

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
			ExpectationStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		body, _ := json.Marshal(tc.ReqBody)
		reqBody := bytes.NewReader(body)
		req, _ := http.NewRequest(http.MethodPost, "/transfer", reqBody)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		serverTest.router.ServeHTTP(rr, req)

		if rr.Code != tc.ExpectationStatusCode {
			t.Fatalf("failed %s wrong response code, want %d got %d", tc.Name, tc.ExpectationStatusCode, rr.Code)
		}
	}
}
