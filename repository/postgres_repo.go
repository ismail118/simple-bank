package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
	"log"
	"time"
)

type PostgresRepository struct {
	db DBTX
}

func NewPostgresRepo(db DBTX) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) WithTx(tx *sql.Tx) Repository {
	return &PostgresRepository{
		db: tx,
	}
}

// InsertAccount insert new account to database and return newID and error if exist
func (r *PostgresRepository) InsertAccount(ctx context.Context, arg models.Account) (int64, error) {
	var newID int64
	query := `
	insert into accounts (owner, balance, currency, created_at)
	values ($1, $2, $3, $4)
	returning id
`
	row := r.db.QueryRowContext(ctx, query,
		arg.Owner,
		arg.Balance,
		arg.Currency,
		time.Now(),
	)

	err := row.Scan(
		&newID,
	)
	if err != nil {
		return 0, err
	}

	err = row.Err()
	if err != nil {
		return newID, err
	}

	return newID, nil
}

// GetAccountByID return account from given id or empty account if not found and error if exist
func (r *PostgresRepository) GetAccountByID(ctx context.Context, id int64) (models.Account, error) {
	query := `
	select id, owner, balance, currency, created_at from accounts
	where id = $1
`
	var a models.Account

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&a.ID,
		&a.Owner,
		&a.Balance,
		&a.Currency,
		&a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("account with id: %d not found in database", id)
			return a, nil
		}
		return a, err
	}

	err = row.Err()
	if err != nil {
		return a, err
	}

	return a, nil
}

// GetListAccounts return list accounts from database and error if exist
func (r *PostgresRepository) GetListAccounts(ctx context.Context, limit, offset int) ([]*models.Account, error) {
	query := `
	select id, owner, balance, currency, created_at from accounts 
	order by id
	limit $1
	offset $2
`
	var items []*models.Account

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Account
		err = rows.Scan(
			&a.ID,
			&a.Owner,
			&a.Balance,
			&a.Currency,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &a)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateAccountBalanceByID update account balance from given id and return error if exist
func (r *PostgresRepository) UpdateAccountBalanceByID(ctx context.Context, balance, id int64) error {
	query := `
	update accounts set balance = $1
	where id = $2
`
	_, err := r.db.ExecContext(ctx, query, balance, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccount delete account from given id return error if exist
func (r *PostgresRepository) DeleteAccount(ctx context.Context, id int64) error {
	query := `
	delete from accounts where id = $1
`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// InsertEntry insert new entry to database and return newID and error if exist
func (r *PostgresRepository) InsertEntry(ctx context.Context, arg models.Entry) (int64, error) {
	query := `
	insert into entries (account_id, amount, created_at) 
	values ($1, $2, $3)
	returning id
`
	var newID int64
	row := r.db.QueryRowContext(ctx, query,
		arg.AccountID,
		arg.Amount,
		time.Now(),
	)

	err := row.Scan(&newID)
	if err != nil {
		return 0, err
	}

	err = row.Err()
	if err != nil {
		return newID, err
	}

	return newID, nil
}

// GetEntryByID return entry by given id or empty entry if not found and error if exist
func (r *PostgresRepository) GetEntryByID(ctx context.Context, id int64) (models.Entry, error) {
	query := `
	select id, account_id, amount, created_at from entries
	where id = $1
`
	var a models.Entry

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&a.ID,
		&a.AccountID,
		&a.Amount,
		&a.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("entries with id:%d not found", id)
			return a, nil
		}
		return a, err
	}

	err = row.Err()
	if err != nil {
		return a, err
	}

	return a, nil
}

// GetListEntries return list entry from given account_id and error if exist
func (r *PostgresRepository) GetListEntries(ctx context.Context, accountID int64, limit, offset int) ([]*models.Entry, error) {
	query := `
	select id, account_id, amount, created_at from entries
	where account_id = $1
	order by id
	limit $2
	offset $3
`
	var items []*models.Entry

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Entry
		err = rows.Scan(
			&a.ID,
			&a.AccountID,
			&a.Amount,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &a)
	}

	err = rows.Err()
	if err != nil {
		return items, err
	}

	return items, nil
}

// InsertTransfer insert new transfer to database and return newID and error if exist
func (r *PostgresRepository) InsertTransfer(ctx context.Context, arg models.Transfer) (int64, error) {
	query := `
	insert into transfers (from_account_id, to_account_id, amount, created_at) 
	values ($1, $2, $3, $4)
	returning id
`
	var newID int64

	row := r.db.QueryRowContext(ctx, query,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
		time.Now(),
	)
	err := row.Scan(&newID)
	if err != nil {
		return 0, err
	}

	err = row.Err()
	if err != nil {
		return newID, err
	}

	return newID, nil
}

// GetTransferByID return transfers from given id or empty transfers if not found and error if exist
func (r *PostgresRepository) GetTransferByID(ctx context.Context, id int64) (models.Transfer, error) {
	query := `
	select id, from_account_id, to_account_id, amount, created_at from transfers
	where id = $1
`
	var a models.Transfer

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&a.ID,
		&a.FromAccountID,
		&a.ToAccountID,
		&a.Amount,
		&a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("transfers with id:%d not found", id)
			return a, nil
		}
		return a, err
	}

	err = row.Err()
	if err != nil {
		return a, err
	}

	return a, nil
}

// GetListTransfers return list transfers from given from_account_id or to_account_id and error if exist
func (r *PostgresRepository) GetListTransfers(ctx context.Context, fromAccountID, toAccountID int64, limit, offset int) ([]*models.Transfer, error) {
	query := `
	select id, from_account_id, to_account_id, amount, created_at from transfers 
	where from_account_id = $1 or to_account_id = $2
	order by id
	limit $3
	offset $4
`
	var items []*models.Transfer

	rows, err := r.db.QueryContext(ctx, query, fromAccountID, toAccountID, limit, offset)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Transfer
		err = rows.Scan(
			&a.ID,
			&a.FromAccountID,
			&a.ToAccountID,
			&a.Amount,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &a)
	}

	err = rows.Err()
	if err != nil {
		return items, err
	}

	return items, nil
}

// GetAccountByIdForUpdate return account from given id or empty account if not found and error if exist
func (r *PostgresRepository) GetAccountByIdForUpdate(ctx context.Context, id int64) (models.Account, error) {
	query := `
	select id, owner, balance, currency, created_at from accounts
	where id = $1 
	for no key update;
`
	var a models.Account

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&a.ID,
		&a.Owner,
		&a.Balance,
		&a.Currency,
		&a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("account with id: %d not found in database", id)
			return a, nil
		}
		return a, err
	}

	err = row.Err()
	if err != nil {
		return a, err
	}

	return a, nil
}

// AddAccountBalanceByID increase or decrease account balance from given id and return accounts and error if exist.
// if amount argument is positive number it will increase the balance and otherwise if negative number it will decrease the balance
func (r *PostgresRepository) AddAccountBalanceByID(ctx context.Context, amount, id int64) (models.Account, error) {
	query := `
	update accounts set balance = balance + $1
	where id = $2
	returning id, owner, balance, currency, created_at
`
	var a models.Account

	row := r.db.QueryRowContext(ctx, query, amount, id)
	err := row.Scan(
		&a.ID,
		&a.Owner,
		&a.Balance,
		&a.Currency,
		&a.CreatedAt,
	)
	if err != nil {
		return a, err
	}

	return a, nil
}
