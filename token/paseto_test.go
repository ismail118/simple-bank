package token

import (
	"errors"
	"github.com/ismail118/simple-bank/util"
	paseto2 "github.com/o1egl/paseto"
	"testing"
	"time"
)

func Test_PasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	if err != nil {
		t.Fatalf("failed create JWTMaker err:%s", err)
	}

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now().Add(time.Second)
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	if err != nil {
		t.Fatalf("failed create Token err:%s", err)
	}
	if token == "" {
		t.Fatalf("failed token is empty")
	}
	if payload == nil {
		t.Fatalf("failed payload is empty")
	}

	payload, err = maker.VerifyToken(token)
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

func Test_ExpiredPaseto(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	if err != nil {
		t.Fatalf("failed create JWTMaker err:%s", err)
	}

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	if err != nil {
		t.Fatalf("failed create Token err:%s", err)
	}
	if token == "" {
		t.Fatalf("failed token is empty")
	}
	if payload == nil {
		t.Fatalf("failed payload is empty")
	}

	payload, err = maker.VerifyToken(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("failed wrong type error want %s got %s", ErrExpiredToken, err)
	}
	if payload != nil {
		t.Fatalf("failed payload should be empty")
	}
}

func Test_InvalidPaseto(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	if err != nil {
		t.Fatalf("failed NewPayload err:%s", err)
	}

	paseto := paseto2.NewV2()
	token, err := paseto.Encrypt([]byte(util.RandomString(32)), payload, nil)
	if err != nil {
		t.Fatalf("failed encrypt error:%s", err)
	}

	maker, err := NewPasetoMaker(util.RandomString(32))
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
