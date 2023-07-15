package repository

import (
	"context"
	"github.com/ismail118/simple-bank/models"
)

type CreateUserTxResult struct {
	User models.Users
}

func (s *SQLStore) CreateUserTx(ctx context.Context, arg models.Users, afterCreate func(user models.Users) error) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(repo Repository) error {
		err := repo.InsertUsers(ctx, arg)
		if err != nil {
			return err
		}

		err = afterCreate(arg)
		if err != nil {
			return err
		}

		result.User = arg

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
