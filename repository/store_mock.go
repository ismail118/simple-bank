package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
)

type SQLStoreMock struct {
	Repository
	db *sql.DB
}

func NewStoreMock(db *sql.DB) Store {
	return &SQLStoreMock{
		db:         db,
		Repository: NewPostgresRepo(db),
	}
}

func (s *SQLStoreMock) execTx(ctx context.Context, fn func(Repository) error) error {
	return nil
}

func (s *SQLStoreMock) TransferTx(ctx context.Context, arg models.Transfer) (TransferTxResult, error) {
	var result TransferTxResult
	return result, nil
}
