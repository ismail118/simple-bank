package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_authorizationMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "Accepted",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken(util.RandomOwner(), time.Minute)
				if err != nil {
					t.Fatalf("failed crate token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Code != http.StatusAccepted {
					t.Fatalf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
				}
			},
		},
		{
			name: "Auth-not-provided",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken(util.RandomOwner(), time.Minute)
				if err != nil {
					t.Fatalf("failed crate token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				r.Header.Set("WRONG-HEADER", fmt.Sprintf("Bearer %s", token))
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Code != http.StatusUnauthorized {
					t.Fatalf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
				}

				var res map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&res)
				if err != nil {
					t.Fatalf("failed decoded response")
				}

				error, ok := res["error"].(string)
				if !ok {
					t.Fatalf("failed wrong response body")
				}

				var expectedErr string = "authorization header is not provided"
				if error != expectedErr {
					t.Fatalf("faield wrong response body want %s got %s", expectedErr, error)
				}
			},
		},
		{
			name: "Auth-invalid-format",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken(util.RandomOwner(), time.Minute)
				if err != nil {
					t.Fatalf("failed crate token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("INVALID-FORMAT%s", token))
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Code != http.StatusUnauthorized {
					t.Fatalf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
				}

				var res map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&res)
				if err != nil {
					t.Fatalf("failed decoded response")
				}

				error, ok := res["error"].(string)
				if !ok {
					t.Fatalf("failed wrong response body")
				}

				var expectedErr string = "invalid authorization header format"
				if error != expectedErr {
					t.Fatalf("faield wrong response body want %s got %s", expectedErr, error)
				}
			},
		},
		{
			name: "Auth-unsupported-type",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				token, payload, err := tokenMaker.CreateToken(util.RandomOwner(), time.Minute)
				if err != nil {
					t.Fatalf("failed crate token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("UNSUPPORTED %s", token))
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Code != http.StatusUnauthorized {
					t.Fatalf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
				}

				var res map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&res)
				if err != nil {
					t.Fatalf("failed decoded response")
				}

				error, ok := res["error"].(string)
				if !ok {
					t.Fatalf("failed wrong response body")
				}

				var expectedErr string = "unsupported authorization type unsupported"
				if error != expectedErr {
					t.Fatalf("faield wrong response body want %s got %s", expectedErr, error)
				}
			},
		},
		{
			name: "Auth-invalid-token",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				tokenMakerOther, err := token.NewPasetoMaker(util.RandomString(32))
				if err != nil {
					t.Fatalf("failed setup NewPasetoMakerOther error:%s", err)
				}

				token, payload, err := tokenMakerOther.CreateToken(util.RandomOwner(), time.Minute)
				if err != nil {
					t.Fatalf("failed crate token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				r.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				if rr.Code != http.StatusUnauthorized {
					t.Fatalf("failed wrong response code, want %d got %d", http.StatusAccepted, rr.Code)
				}
			},
		},
	}

	router := gin.New()
	tokenMaker := serverTest.tokenMaker

	router.Use(authMiddleware(tokenMaker))
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, "Ok")
	})

	for _, tc := range testCases {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		// setup
		tc.setupAuth(t, req, tokenMaker)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		//check response
		tc.checkResponse(t, rr)
	}

}
