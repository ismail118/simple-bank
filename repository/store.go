package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ismail118/simple-bank/models"
)

type Store interface {
	Repository
	execTx(ctx context.Context, fn func(Repository) error) error
	TransferTx(ctx context.Context, arg models.Transfer) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg models.Users, afterCreate func(user models.Users) error) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, id int64, secretCode string) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	Repository
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:         db,
		Repository: NewPostgresRepo(db),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(Repository) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repo := NewPostgresRepo(tx)
	err = fn(repo)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx error:%s rollback error: %s", err, rbErr)
		}
	}

	return tx.Commit()
}
