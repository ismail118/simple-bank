package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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

func (r *PostgresRepository) GetAccountByOwnerAndCurrency(ctx context.Context, owner, currency string) (models.Account, error) {
	query := `
	select id, owner, balance, currency, created_at from accounts
	where owner = $1 and currency = $2
`
	var a models.Account

	row := r.db.QueryRowContext(ctx, query, owner, currency)
	err := row.Scan(
		&a.ID,
		&a.Owner,
		&a.Balance,
		&a.Currency,
		&a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("account with owner %s and currency %s not found in database", owner, currency)
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
func (r *PostgresRepository) GetListAccounts(ctx context.Context, owner string, limit, offset int) ([]*models.Account, error) {
	query := `
	select id, owner, balance, currency, created_at 
	from accounts 
	where owner = $1 
	order by id
	limit $2
	offset $3
`
	items := []*models.Account{}

	rows, err := r.db.QueryContext(ctx, query, owner, limit, offset)
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

// UpdateAccount update account balance from given id and return error if exist
func (r *PostgresRepository) UpdateAccount(ctx context.Context, arg models.Account) error {
	query := `
	update accounts set balance = $1, currency = $2
	where id = $3
`
	_, err := r.db.ExecContext(ctx, query, arg.Balance, arg.Currency, arg.ID)
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
	items := []*models.Entry{}

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
	items := []*models.Transfer{}

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

func (r *PostgresRepository) InsertUsers(ctx context.Context, arg models.Users) error {
	query := `
	insert into users (username, hashed_password, full_name, email, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6)
`

	_, err := r.db.ExecContext(ctx, query,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetUsersByUsername(ctx context.Context, username string) (models.Users, error) {
	query := `
	select username, hashed_password, full_name, email, created_at, updated_at, is_email_verify from users
	where username = $1
`
	var a models.Users
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&a.Username,
		&a.HashedPassword,
		&a.FullName,
		&a.Email,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.IsEmailVerify,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("users with username %s not found", username)
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

func (r *PostgresRepository) GetUsersByEmail(ctx context.Context, email string) (models.Users, error) {
	query := `
	select username, hashed_password, full_name, email, created_at, updated_at, is_email_verify from users
	where email = $1
`
	var a models.Users
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&a.Username,
		&a.HashedPassword,
		&a.FullName,
		&a.Email,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.IsEmailVerify,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("users with email %s not found", email)
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

func (r *PostgresRepository) GetListUsers(ctx context.Context, limit, offset int) ([]*models.Users, error) {
	query := `
	select username, hashed_password, full_name, email, created_at, updated_at, is_email_verify from users limit $1 offset $2
`
	var items []*models.Users
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a models.Users
		err := rows.Scan(
			&a.Username,
			&a.HashedPassword,
			&a.FullName,
			&a.Email,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.IsEmailVerify,
		)
		if err != nil {
			return items, err
		}
		items = append(items, &a)
	}

	err = rows.Err()
	if err != nil {
		return items, err
	}

	return items, nil
}

type UpdateUserParam struct {
	Username        string         `json:"username"`
	HashedPassword  sql.NullString `json:"hashed_password"`
	FullName        sql.NullString `json:"full_name"`
	Email           sql.NullString `json:"email"`
	IsEmailVerified sql.NullBool   `json:"is_email_verified"`
}

func (r *PostgresRepository) UpdateUsers(ctx context.Context, arg UpdateUserParam) error {
	query := `
	update users 
	set 
	    hashed_password = COALESCE($1, hashed_password),
	    full_name = COALESCE($2, full_name),
	    email = COALESCE($3, email),
	    is_email_verify = COALESCE($4, is_email_verify),
	    updated_at = $5
	where username = $6
`
	_, err := r.db.ExecContext(ctx, query,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
		time.Now(),
		arg.Username,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) DeleteUsers(ctx context.Context, username string) error {
	query := `
	delete from users where username = $1
`
	_, err := r.db.ExecContext(ctx, query, username)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) InsertSessions(ctx context.Context, arg models.Sessions) error {
	query := `
	insert into sessions (id, username, refresh_token, user_agent, client_ip, is_blocked, expired_at, created_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8) 
`
	_, err := r.db.ExecContext(ctx, query,
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		false,
		arg.ExpiredAt,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetSessionsByID(ctx context.Context, id uuid.UUID) (models.Sessions, error) {
	query := `
	select id, username, refresh_token, user_agent, client_ip, is_blocked, expired_at, created_at
	from sessions
	where id = $1
`
	var s models.Sessions
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&s.ID,
		&s.Username,
		&s.RefreshToken,
		&s.UserAgent,
		&s.ClientIp,
		&s.IsBlocked,
		&s.ExpiredAt,
		&s.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, nil
		}
		return s, err
	}

	err = row.Err()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (r *PostgresRepository) InsertVerifyEmail(ctx context.Context, arg models.VerifyEmail) (int64, error) {
	query := `
	insert into verify_email (username, email, secret_code, is_used, created_at, expired_at)
	values ($1, $2, $3, $4, $5, $6)
	returning id
`
	var id int64
	row := r.db.QueryRowContext(ctx, query,
		arg.Username,
		arg.Email,
		arg.SecretCode,
		false,
		time.Now(),
		time.Now().Add(15*time.Minute),
	)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	err = row.Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresRepository) GetVerifyEmailByID(ctx context.Context, id int64) (models.VerifyEmail, error) {
	query := `
	select id, username, email, secret_code, is_used, created_at, expired_at
	from verify_email
	where id = $1
`
	var a models.VerifyEmail
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&a.ID,
		&a.Username,
		&a.Email,
		&a.SecretCode,
		&a.IsUsed,
		&a.CreatedAt,
		&a.ExpiredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
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

// UpdateVerifyEmailIsUsed update table verify_email field is_used to true, return error if nothing row effected
func (r *PostgresRepository) UpdateVerifyEmailIsUsed(ctx context.Context, id int64, secretCode string) error {
	query := `
	update verify_email set is_used = true
	where id = $1
	  AND secret_code = $2
	  AND is_used = false
	  AND expired_at > now()
`
	res, err := r.db.ExecContext(ctx, query, id, secretCode)
	if err != nil {
		return err
	}

	effected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if effected == 0 {
		return fmt.Errorf("nothing row effected")
	}
	return nil
}
