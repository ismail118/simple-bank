package repository

import (
	"context"
	"github.com/ismail118/simple-bank/models"
)

type TransferTxResult struct {
	Transfer    models.Transfer `json:"transfer"`
	FromAccount models.Account  `json:"from_account"`
	ToAccount   models.Account  `json:"to_account"`
	FromEntry   models.Entry    `json:"from_entry"`
	ToEntry     models.Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg models.Transfer) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(r Repository) error {
		newID, err := s.InsertTransfer(ctx, arg)
		if err != nil {
			return err
		}
		arg.ID = newID

		result.Transfer = arg

		fEntry := models.Entry{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}
		newID, err = s.InsertEntry(ctx, fEntry)
		if err != nil {
			return err
		}
		fEntry.ID = newID

		result.FromEntry = fEntry

		tEntry := models.Entry{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}
		newID, err = s.InsertEntry(ctx, tEntry)
		if err != nil {
			return err
		}
		tEntry.ID = newID

		result.ToEntry = tEntry

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = s.AddAccountBalanceByID(context.Background(), -arg.Amount, arg.FromAccountID)
			if err != nil {
				return err
			}
			result.ToAccount, err = s.AddAccountBalanceByID(context.Background(), arg.Amount, arg.ToAccountID)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = s.AddAccountBalanceByID(context.Background(), arg.Amount, arg.ToAccountID)
			if err != nil {
				return err
			}
			result.FromAccount, err = s.AddAccountBalanceByID(context.Background(), -arg.Amount, arg.FromAccountID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
