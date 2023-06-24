package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ismail118/simple-bank/util"
	"testing"
	"time"
)

func Test_JWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatalf("failed create JWTMaker err:%s", err)
	}

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now().Add(time.Second)
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	if err != nil {
		t.Fatalf("failed create Token err:%s", err)
	}
	if token == "" {
		t.Fatalf("failed token is empty")
	}

	payload, err := maker.VerifyToken(token)
	if err != nil {
		t.Fatalf("failed verify Token err:%s", err)
	}
	if payload == nil {
		t.Fatalf("failed payload is empty")
	}

	if payload.ID.String() == "" {
		t.Fatalf("failed payload id empty")
	}
	if payload.Username != username {
		t.Fatalf("failed mismatch username")
	}
	if !payload.IssuedAt.Before(issuedAt) {
		t.Fatalf("failed mismatch issued_at")
	}
	if !payload.ExpiredAt.Before(expiredAt) {
		t.Fatalf("failed mismatch expired_at")
	}
}

func Test_ExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatalf("failed create JWTMaker err:%s", err)
	}

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	if err != nil {
		t.Fatalf("failed create Token err:%s", err)
	}
	if token == "" {
		t.Fatalf("failed token is empty")
	}

	payload, err := maker.VerifyToken(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("failed wrong type error want %s got %s", ErrExpiredToken, err)
	}
	if payload != nil {
		t.Fatalf("failed payload should be empty")
	}
}

func Test_InvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	if err != nil {
		t.Fatalf("failed NewPayload err:%s", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed SignedString error:%s", err)
	}

	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatalf("failed create JWTMaker err:%s", err)
	}

	payload, err = maker.VerifyToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("failed wrong type error want %s got %s", ErrInvalidToken, err)
	}
	if payload != nil {
		t.Fatalf("failed payload should be empty")
	}
}
