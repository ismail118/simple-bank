package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
)

type Repository interface {
	InsertAccount(ctx context.Context, arg models.Account) (int64, error)
	GetAccountByID(ctx context.Context, id int64) (models.Account, error)
	GetListAccounts(ctx context.Context, owner string, limit, offset int) ([]*models.Account, error)
	UpdateAccount(ctx context.Context, arg models.Account) error
	DeleteAccount(ctx context.Context, id int64) error
	InsertEntry(ctx context.Context, arg models.Entry) (int64, error)
	GetEntryByID(ctx context.Context, id int64) (models.Entry, error)
	GetListEntries(ctx context.Context, accountID int64, limit, offset int) ([]*models.Entry, error)
	InsertTransfer(ctx context.Context, arg models.Transfer) (int64, error)
	GetTransferByID(ctx context.Context, id int64) (models.Transfer, error)
	GetListTransfers(ctx context.Context, fromAccountID, toAccountID int64, limit, offset int) ([]*models.Transfer, error)
	GetAccountByIdForUpdate(ctx context.Context, id int64) (models.Account, error)
	AddAccountBalanceByID(ctx context.Context, amount, id int64) (models.Account, error)
	InsertUsers(ctx context.Context, arg models.Users) error
	GetUsersByUsername(ctx context.Context, username string) (models.Users, error)
	GetListUsers(ctx context.Context, limit, offset int) ([]*models.Users, error)
	UpdateUsers(ctx context.Context, arg models.Users) error
	DeleteUsers(ctx context.Context, username string) error
	GetAccountByOwnerAndCurrency(ctx context.Context, owner, currency string) (models.Account, error)
	GetUsersByEmail(ctx context.Context, email string) (models.Users, error)
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
