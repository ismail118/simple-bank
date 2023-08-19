package api

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ismail118/simple-bank/models"
	pb "github.com/ismail118/simple-bank/proto"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/util"
	"github.com/ismail118/simple-bank/worker"
	"github.com/stretchr/testify/assert"
)


func Test_Grpc_CreateUser(t *testing.T) {
	user, _, err := util.RandomUser()
	assert.NoError(t, err)

	testCases := []struct {
		name string
		req *pb.CreateUserRequest
		buildStub func(t *testing.T, taskDistributor *worker.MockTaskDistributor, store *repository.MockStore, req *pb.CreateUserRequest)
		checkResponse func(t *testing.T, resp *pb.CreateUserResponse, err error)
	} {
		{
			name: "ok",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email: user.Email,
				Password: util.RandomString(10),
			},
			buildStub: func(t *testing.T, taskDistributor *worker.MockTaskDistributor, store *repository.MockStore, req *pb.CreateUserRequest) {
				hashedPassword, err := util.HashedPassword(req.GetPassword())
				assert.NoError(t, err)
				
				user.HashedPassword = hashedPassword

				arg := models.Users{
					Username:       req.GetUsername(),
					FullName:       req.GetFullName(),
					Email:          req.GetEmail(),
				}

				// taskPayload := &worker.PayloadSendVerifyEmail{Username: req.GetUsername()}

				store.EXPECT().
				CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, req.GetPassword(), user), gomock.Any()).
				Times(1).
				Return(repository.CreateUserTxResult{
					User: user,
				}, nil)

				// taskDistributor.EXPECT().
				// DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Eq(taskPayload), gomock.Any()).
				// Times(1).
				// Return(nil)
			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				createdUser := resp.GetUser()
				assert.Equal(t, user.Username, createdUser.Username)
				assert.Equal(t, user.FullName, createdUser.FullName)
				assert.Equal(t, user.Email, createdUser.Email)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := repository.NewMockStore(ctrl)
		repo := repository.NewMockRepository(ctrl)
		taskDistributor := worker.NewMockTaskDistributor(ctrl)

		tc.buildStub(t, taskDistributor, store, tc.req)

		srv := NewGrpcServer(store, repo, tokenMakerTest, &util.Config{AccessTokenDuration: time.Minute, RefreshTokenDuration: time.Minute}, taskDistributor)

		res, err := srv.CreateUser(context.Background(), tc.req)
		tc.checkResponse(t, res, err)
	}
}
