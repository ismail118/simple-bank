package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
)

type Repository interface {
	InsertAccount(ctx context.Context, arg models.Account) (int64, error)
	GetAccountByID(ctx context.Context, id int64) (models.Account, error)
	GetListAccounts(ctx context.Context, limit, offset int) ([]*models.Account, error)
	UpdateAccountBalanceByID(ctx context.Context, balance, id int64) error
	DeleteAccount(ctx context.Context, id int64) error
	InsertEntry(ctx context.Context, arg models.Entry) (int64, error)
	GetEntryByID(ctx context.Context, id int64) (models.Entry, error)
	GetListEntries(ctx context.Context, accountID int64, limit, offset int) ([]*models.Entry, error)
	InsertTransfer(ctx context.Context, arg models.Transfer) (int64, error)
	GetTransferByID(ctx context.Context, id int64) (models.Transfer, error)
	GetListTransfers(ctx context.Context, fromAccountID, toAccountID int64, limit, offset int) ([]*models.Transfer, error)
	GetAccountByIdForUpdate(ctx context.Context, id int64) (models.Account, error)
	AddAccountBalanceByID(ctx context.Context, amount, id int64) (models.Account, error)
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
