package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/ismail118/simple-bank/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	UpdateUsers(ctx context.Context, arg UpdateUserParam) error
	DeleteUsers(ctx context.Context, username string) error
	GetAccountByOwnerAndCurrency(ctx context.Context, owner, currency string) (models.Account, error)
	GetUsersByEmail(ctx context.Context, email string) (models.Users, error)
	InsertSessions(ctx context.Context, arg models.Sessions) error
	GetSessionsByID(ctx context.Context, id uuid.UUID) (models.Sessions, error)
	InsertVerifyEmail(ctx context.Context, arg models.VerifyEmail) (int64, error)
	GetVerifyEmailByID(ctx context.Context, id int64) (models.VerifyEmail, error)
	UpdateVerifyEmailIsUsed(ctx context.Context, id int64, secretCode string) error
}

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Manager struct {
    FullName       string
    Position       string
    Age            int32
    YearsInCompany int32
}