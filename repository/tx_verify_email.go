package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
)

type VerifyEmailTxResult struct {
	VerifyEmail models.VerifyEmail
	User        models.Users
}

func (s *SQLStore) VerifyEmailTx(ctx context.Context, id int64, secretCode string) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := s.execTx(ctx, func(repo Repository) error {
		err := repo.UpdateVerifyEmailIsUsed(ctx, id, secretCode)
		if err != nil {
			return err
		}

		verifyEmail, err := repo.GetVerifyEmailByID(ctx, id)
		if err != nil {
			return err
		}
		result.VerifyEmail = verifyEmail

		err = repo.UpdateUsers(ctx, UpdateUserParam{
			Username: verifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		result.User, err = repo.GetUsersByUsername(ctx, verifyEmail.Username)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
