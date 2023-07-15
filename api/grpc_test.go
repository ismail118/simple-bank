package api

import (
	"context"
	"fmt"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func Test_Grpc_CreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		isNilResponse bool
		isError       bool
	}{
		{
			name: "ok",
			req: &pb.CreateUserRequest{
				Username: "user",
				FullName: util.RandomString(6),
				Email:    "notexists@gmail.com",
				Password: util.RandomString(12),
			},
			isNilResponse: false,
			isError:       false,
		},
		{
			name: "invalid-request",
			req: &pb.CreateUserRequest{
				Username: "",
				FullName: util.RandomString(6),
				Email:    "notexists@gmail.com",
				Password: util.RandomString(12),
			},
			isNilResponse: true,
			isError:       true,
		},
		{
			name: "user-already-exist",
			req: &pb.CreateUserRequest{
				Username: "ismail",
				FullName: util.RandomString(6),
				Email:    "notexists@gmail.com",
				Password: util.RandomString(12),
			},
			isNilResponse: true,
			isError:       true,
		},
		{
			name: "email-already-exist",
			req: &pb.CreateUserRequest{
				Username: "user",
				FullName: util.RandomString(6),
				Email:    "exist@gmail.com",
				Password: util.RandomString(12),
			},
			isNilResponse: true,
			isError:       true,
		},
		{
			name: "error-send-verify-email",
			req: &pb.CreateUserRequest{
				Username: "user2",
				FullName: util.RandomString(6),
				Email:    "notexists@gmail.com",
				Password: util.RandomString(12),
			},
			isNilResponse: true,
			isError:       true,
		},
	}

	for _, tc := range testCases {
		res, err := grpcServerTest.CreateUser(context.Background(), tc.req)

		if tc.isNilResponse {
			if res != nil {
				t.Errorf("failed %s res should be nil: %v", tc.name, res)
			}
		} else {
			if res == nil {
				t.Errorf("failed %s res should be not nil", tc.name)
			}
		}

		if tc.isError {
			if err == nil {
				t.Fatalf("failed %s error should not be nil", tc.name)
			}
		} else {
			if err != nil {
				t.Fatalf("failed %s error should be nil: %s", tc.name, err)
			}
		}
	}
}

func Test_Grpc_UpdateUser(t *testing.T) {
	testCases := []struct {
		name       string
		req        *pb.UpdateUserRequest
		setupToken func(username string, tokenMaker token.Maker) context.Context
		isError    bool
	}{
		{
			name: "ok",
			req: &pb.UpdateUserRequest{
				Username: "user3",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("notexists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				token, payload, err := tokenMaker.CreateToken(username, time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				md := metadata.Pairs(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
				return metadata.NewIncomingContext(context.Background(), md)
			},
			isError: false,
		},
		{
			name: "unauthorize",
			req: &pb.UpdateUserRequest{
				Username: "user3",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("notexists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			isError: true,
		},
		{
			name: "invalid-argument",
			req: &pb.UpdateUserRequest{
				Username: "1",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("notexists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				token, payload, err := tokenMaker.CreateToken(username, time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				md := metadata.Pairs(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
				return metadata.NewIncomingContext(context.Background(), md)
			},
			isError: true,
		},
		{
			name: "username-not-found",
			req: &pb.UpdateUserRequest{
				Username: "user",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("exists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				token, payload, err := tokenMaker.CreateToken(username, time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				md := metadata.Pairs(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
				return metadata.NewIncomingContext(context.Background(), md)
			},
			isError: true,
		},
		{
			name: "email-already-used",
			req: &pb.UpdateUserRequest{
				Username: "user3",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("exists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				token, payload, err := tokenMaker.CreateToken(username, time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				md := metadata.Pairs(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
				return metadata.NewIncomingContext(context.Background(), md)
			},
			isError: true,
		},
		{
			name: "unauthorized",
			req: &pb.UpdateUserRequest{
				Username: "user3",
				FullName: toPointerString(util.RandomOwner()),
				Email:    toPointerString("notexists@gmail.com"),
				Password: nil,
			},
			setupToken: func(username string, tokenMaker token.Maker) context.Context {
				token, payload, err := tokenMaker.CreateToken("notme", time.Minute)
				if err != nil {
					t.Fatalf("failed create token error:%s", err)
				}
				if payload == nil {
					t.Fatalf("failed payload is empty")
				}

				md := metadata.Pairs(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, token))
				return metadata.NewIncomingContext(context.Background(), md)
			},
			isError: true,
		},
	}

	tokenMaker := grpcServerTest.tokenMaker

	for _, tc := range testCases {
		ctx := tc.setupToken(tc.req.Username, tokenMaker)
		_, err := grpcServerTest.UpdateUser(ctx, tc.req)
		if tc.isError {
			if err == nil {
				t.Fatalf("failed %s error should be nil", tc.name)
			}
		} else {
			if err != nil {
				t.Fatalf("failed %s error should be nil:%s", tc.name, err)
			}
		}
	}
}

func Test_Grpc_Login(t *testing.T) {
	testCases := []struct {
		name    string
		req     *pb.LoginRequest
		isError bool
	}{
		{
			name: "ok",
			req: &pb.LoginRequest{
				Username: "user3",
				Password: "some password",
			},
			isError: false,
		},
		{
			name: "user-not-found",
			req: &pb.LoginRequest{
				Username: "user",
				Password: "some password",
			},
			isError: true,
		},
		{
			name: "invalid-argument",
			req: &pb.LoginRequest{
				Username: "",
				Password: "",
			},
			isError: true,
		},
		{
			name: "wrong-password",
			req: &pb.LoginRequest{
				Username: "user3",
				Password: "wrong password",
			},
			isError: true,
		},
	}

	for _, tc := range testCases {
		_, err := grpcServerTest.Login(context.Background(), tc.req)
		if tc.isError {
			if err == nil {
				t.Fatalf("failed %s error should be nil", tc.name)
			}
		} else {
			if err != nil {
				t.Fatalf("failed %s error should be nil:%s", tc.name, err)
			}
		}
	}
}

func Test_Grpc_VerifyEmail(t *testing.T) {
	testCases := []struct {
		name    string
		req     *pb.VerifyEmailRequest
		isError bool
	}{
		{
			name: "ok",
			req: &pb.VerifyEmailRequest{
				Id:         1,
				SecretCode: util.RandomString(32),
			},
			isError: false,
		},
	}

	for _, tc := range testCases {
		_, err := grpcServerTest.VerifyEmail(context.Background(), tc.req)
		if tc.isError {
			if err == nil {
				t.Fatalf("failed %s error should be not nil", tc.name)
			}
		} else {
			if err != nil {
				t.Fatalf("failed %s error should be nil: %s", tc.name, err)
			}
		}
	}
}

func toPointerString(s string) *string {
	return &s
}
