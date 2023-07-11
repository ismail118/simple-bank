package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ismail118/simple-bank/models"
	"time"
)

type PostgresRepositoryMock struct {
	db DBTX
}

func NewPostgresRepoMock(db DBTX) Repository {
	return &PostgresRepositoryMock{
		db: db,
	}
}

func (r *PostgresRepositoryMock) WithTx(tx *sql.Tx) Repository {
	return &PostgresRepositoryMock{
		db: tx,
	}
}

// InsertAccount insert new account to database and return newID and error if exist
func (r *PostgresRepositoryMock) InsertAccount(ctx context.Context, arg models.Account) (int64, error) {
	var newID int64
	if arg.Owner == "test-error-db" {
		return newID, sql.ErrConnDone
	}
	return newID, nil
}

// GetAccountByID return account from given id or empty account if not found and error if exist
func (r *PostgresRepositoryMock) GetAccountByID(ctx context.Context, id int64) (models.Account, error) {
	var a models.Account
	if id == 2 || id == 3 {
		a = models.Account{
			ID:        id,
			Owner:     "some-user",
			Balance:   100,
			Currency:  "USD",
			CreatedAt: time.Now(),
		}
	}
	if id > 1000 {
		return a, sql.ErrConnDone
	}
	return a, nil
}

func (r *PostgresRepositoryMock) GetAccountByOwnerAndCurrency(ctx context.Context, username, currency string) (models.Account, error) {
	var a models.Account
	return a, nil
}

// GetListAccounts return list accounts from database and error if exist
func (r *PostgresRepositoryMock) GetListAccounts(ctx context.Context, owner string, limit, offset int) ([]*models.Account, error) {
	items := []*models.Account{}
	if offset > 1000 {
		return items, sql.ErrConnDone
	}
	return items, nil
}

// UpdateAccount update account balance from given id and return error if exist
func (r *PostgresRepositoryMock) UpdateAccount(ctx context.Context, arg models.Account) error {
	if arg.ID == 3 {
		return sql.ErrConnDone
	}
	return nil
}

// DeleteAccount delete account from given id return error if exist
func (r *PostgresRepositoryMock) DeleteAccount(ctx context.Context, id int64) error {
	if id == 3 {
		return sql.ErrConnDone
	}
	return nil
}

// InsertEntry insert new entry to database and return newID and error if exist
func (r *PostgresRepositoryMock) InsertEntry(ctx context.Context, arg models.Entry) (int64, error) {
	var newID int64
	return newID, nil
}

// GetEntryByID return entry by given id or empty entry if not found and error if exist
func (r *PostgresRepositoryMock) GetEntryByID(ctx context.Context, id int64) (models.Entry, error) {
	var a models.Entry
	if id == 2 {
		a = models.Entry{
			ID:        id,
			AccountID: 2,
			Amount:    100,
			CreatedAt: time.Now(),
		}
	}
	if id > 1000 {
		return a, sql.ErrConnDone
	}
	return a, nil
}

// GetListEntries return list entry from given account_id and error if exist
func (r *PostgresRepositoryMock) GetListEntries(ctx context.Context, accountID int64, limit, offset int) ([]*models.Entry, error) {
	items := []*models.Entry{}
	if offset > 1000 {
		return items, sql.ErrConnDone
	}
	return items, nil
}

// InsertTransfer insert new transfer to database and return newID and error if exist
func (r *PostgresRepositoryMock) InsertTransfer(ctx context.Context, arg models.Transfer) (int64, error) {
	var newID int64
	return newID, nil
}

// GetTransferByID return transfers from given id or empty transfers if not found and error if exist
func (r *PostgresRepositoryMock) GetTransferByID(ctx context.Context, id int64) (models.Transfer, error) {
	var a models.Transfer
	if id == 2 {
		a = models.Transfer{
			ID:            id,
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        10,
			CreatedAt:     time.Now(),
		}
	}
	if id > 1000 {
		return a, sql.ErrConnDone
	}
	return a, nil
}

// GetListTransfers return list transfers from given from_account_id or to_account_id and error if exist
func (r *PostgresRepositoryMock) GetListTransfers(ctx context.Context, fromAccountID, toAccountID int64, limit, offset int) ([]*models.Transfer, error) {
	items := []*models.Transfer{}
	if offset > 1000 {
		return items, sql.ErrConnDone
	}
	return items, nil
}

// GetAccountByIdForUpdate return account from given id or empty account if not found and error if exist
func (r *PostgresRepositoryMock) GetAccountByIdForUpdate(ctx context.Context, id int64) (models.Account, error) {
	var a models.Account
	return a, nil
}

// AddAccountBalanceByID increase or decrease account balance from given id and return accounts and error if exist.
// if amount argument is positive number it will increase the balance and otherwise if negative number it will decrease the balance
func (r *PostgresRepositoryMock) AddAccountBalanceByID(ctx context.Context, amount, id int64) (models.Account, error) {
	var a models.Account
	return a, nil
}

func (r *PostgresRepositoryMock) InsertUsers(ctx context.Context, arg models.Users) error {
	return nil
}

func (r *PostgresRepositoryMock) GetUsersByUsername(ctx context.Context, username string) (models.Users, error) {
	var a models.Users = models.Users{
		Username:       username,
		HashedPassword: "some password",
		FullName:       "some name",
		Email:          "some@email.com",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if username == "user" {
		a = models.Users{}
	}
	return a, nil
}

func (r *PostgresRepositoryMock) GetUsersByEmail(ctx context.Context, email string) (models.Users, error) {
	var a models.Users = models.Users{
		Username:       "some-owner",
		HashedPassword: "some password",
		FullName:       "some name",
		Email:          email,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if email == "notexists@gmail.com" {
		a = models.Users{}
	}
	return a, nil
}

func (r *PostgresRepositoryMock) GetListUsers(ctx context.Context, limit, offset int) ([]*models.Users, error) {
	items := []*models.Users{}
	if offset > 1000 {
		return items, sql.ErrConnDone
	}
	return items, nil
}

func (r *PostgresRepositoryMock) UpdateUsers(ctx context.Context, arg UpdateUserParam) error {
	if arg.Username == "error" {
		return sql.ErrConnDone
	}
	return nil
}

func (r *PostgresRepositoryMock) DeleteUsers(ctx context.Context, username string) error {
	return nil
}

func (r *PostgresRepositoryMock) InsertSessions(ctx context.Context, arg models.Sessions) error {
	return nil
}

func (r *PostgresRepositoryMock) GetSessionsByID(ctx context.Context, id uuid.UUID) (models.Sessions, error) {
	var s models.Sessions

	return s, nil
}
