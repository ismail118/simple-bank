package repository

import (
	"context"
	"fmt"
	"github.com/ismail118/simple-bank/models"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dbpool *pgxpool.Pool
}

func NewStore(dbpool *pgxpool.Pool) Store {
	return &SQLStore{
		dbpool:     dbpool,
		Repository: NewPostgresRepo(dbpool),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(Repository) error) error {
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	repo := NewPostgresRepo(tx)
	err = fn(repo)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("tx error:%s rollback error: %s", err, rbErr)
		}
	}

	return tx.Commit(ctx)
}
