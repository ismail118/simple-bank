package repository

import (
	"context"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type SQLStoreMock struct {
	Repository
	dbpool *pgxpool.Pool
}

func NewStoreMock(dbpool *pgxpool.Pool) Store {
	return &SQLStoreMock{
		dbpool:     dbpool,
		Repository: NewPostgresRepo(dbpool),
	}
}

func (s *SQLStoreMock) execTx(ctx context.Context, fn func(Repository) error) error {
	return nil
}

func (s *SQLStoreMock) TransferTx(ctx context.Context, arg models.Transfer) (TransferTxResult, error) {
	var result TransferTxResult
	return result, nil
}

func (s *SQLStoreMock) CreateUserTx(ctx context.Context, arg models.Users, afterCreate func(user models.Users) error) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := afterCreate(arg)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SQLStoreMock) VerifyEmailTx(ctx context.Context, id int64, secretCode string) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	if id == 1 {
		result = VerifyEmailTxResult{
			VerifyEmail: models.VerifyEmail{
				ID:         1,
				Username:   "user3",
				Email:      "user3@email.com",
				SecretCode: util.RandomString(32),
				IsUsed:     true,
				CreatedAt:  time.Now(),
				ExpiredAt:  time.Now(),
			},
			User: models.Users{
				Username:       "user3",
				HashedPassword: "secret",
				FullName:       "user user",
				Email:          "user3@email.com",
				IsEmailVerify:  true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
	}
	return result, nil
}
