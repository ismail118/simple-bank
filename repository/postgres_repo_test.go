package repository

import (
	"context"
	"database/sql"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/util"
	"testing"
	"time"
)

func createRandomAccount(owner string) models.Account {
	return models.Account{
		Owner:    owner,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func TestInsertAccount(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))

	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	dataTest := createRandomAccount(user.Username)

	// test insert
	newID, err := testRepo.InsertAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	// get back the data
	data, err := testRepo.GetAccountByID(context.Background(), newID)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if dataTest.Owner != data.Owner {
		t.Errorf("failed deffrent owner, wan't %s got %s", dataTest.Owner, data.Owner)
	}
	if dataTest.Balance != data.Balance {
		t.Errorf("failed deffrent balance, wan't %d got %d", dataTest.Balance, data.Balance)
	}
	if dataTest.Currency != data.Currency {
		t.Errorf("failed deffrent currency, wan't %s got %s", dataTest.Currency, data.Currency)
	}
}

func TestGetAccountByID(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))

	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	dataTest := createRandomAccount(user.Username)

	// test insert
	newID, err := testRepo.InsertAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	// get back the data
	data, err := testRepo.GetAccountByID(context.Background(), newID)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if dataTest.Owner != data.Owner {
		t.Errorf("failed deffrent owner, wan't %s got %s", dataTest.Owner, data.Owner)
	}
	if dataTest.Balance != data.Balance {
		t.Errorf("failed deffrent balance, wan't %d got %d", dataTest.Balance, data.Balance)
	}
	if dataTest.Currency != data.Currency {
		t.Errorf("failed deffrent currency, wan't %s got %s", dataTest.Currency, data.Currency)
	}
}

func TestGetAccountByIdForUpdate(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))

	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	dataTest := createRandomAccount(user.Username)

	// test insert
	newID, err := testRepo.InsertAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	// get back the data
	data, err := testRepo.GetAccountByIdForUpdate(context.Background(), newID)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if dataTest.Owner != data.Owner {
		t.Errorf("failed deffrent owner, wan't %s got %s", dataTest.Owner, data.Owner)
	}
	if dataTest.Balance != data.Balance {
		t.Errorf("failed deffrent balance, wan't %d got %d", dataTest.Balance, data.Balance)
	}
	if dataTest.Currency != data.Currency {
		t.Errorf("failed deffrent currency, wan't %s got %s", dataTest.Currency, data.Currency)
	}
}

func TestUpdateAccount(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))

	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	dataTest := createRandomAccount(user.Username)

	// test insert
	newID, err := testRepo.InsertAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	// update
	dataTest.ID = newID
	dataTest.Balance = util.RandomBalance()

	err = testRepo.UpdateAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed update account balance error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetAccountByID(context.Background(), newID)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if dataTest.Owner != data.Owner {
		t.Errorf("failed deffrent owner, wan't %s got %s", dataTest.Owner, data.Owner)
	}
	if dataTest.Balance != data.Balance {
		t.Errorf("failed deffrent balance, wan't %d got %d", dataTest.Balance, data.Balance)
	}
	if dataTest.Currency != data.Currency {
		t.Errorf("failed deffrent currency, wan't %s got %s", dataTest.Currency, data.Currency)
	}
}

func TestDeleteAccount(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))

	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	dataTest := createRandomAccount(user.Username)

	// test insert
	newID, err := testRepo.InsertAccount(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	// delete
	err = testRepo.DeleteAccount(context.Background(), newID)
	if err != nil {
		t.Errorf("failed delete account error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetAccountByID(context.Background(), newID)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if data.ID > 0 {
		t.Errorf("failed id must be 0")
	}
}

func TestGetListAccounts(t *testing.T) {
	var lastAccount models.Account

	for i := 0; i < 10; i++ {
		hashedPassword, _ := util.HashedPassword(util.RandomString(6))
		user := models.Users{
			Username:       util.RandomOwner(),
			HashedPassword: hashedPassword,
			FullName:       util.RandomOwner(),
			Email:          util.RandomEmail(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// test insert
		err := testRepo.InsertUsers(context.Background(), user)
		if err != nil {
			t.Errorf("failed insert users error:%s", err)
		}

		dataTest := createRandomAccount(user.Username)

		// test insert
		newID, err := testRepo.InsertAccount(context.Background(), dataTest)
		if err != nil {
			t.Errorf("failed insert account error:%s", err)
		}
		if newID < 1 {
			t.Errorf("failed renturn less than 1 newID:%d", newID)
		}

		lastAccount = dataTest
	}

	listData, err := testRepo.GetListAccounts(context.Background(), lastAccount.Owner, 10, 0)
	if err != nil {
		t.Errorf("failed get list accounts error:%s", err)
	}

	if len(listData) < 1 {
		t.Errorf("failed len data not 10 len:%d", len(listData))
	}
}

func TestInsertEntry(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	acc := createRandomAccount(user.Username)

	// test insert account
	newID, err := testRepo.InsertAccount(context.Background(), acc)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	acc.ID = newID

	// insert entry
	entry := models.Entry{
		AccountID: acc.ID,
		Amount:    acc.Balance,
	}

	newID, err = testRepo.InsertEntry(context.Background(), entry)
	if err != nil {
		t.Errorf("failed insert entry error:%s", err)
	}

	entry.ID = newID

	// get back entry
	entry2, err := testRepo.GetEntryByID(context.Background(), entry.ID)
	if err != nil {
		t.Errorf("failed get entry error:%s", err)
	}

	if entry.ID != entry2.ID {
		t.Errorf("failed deffrent entry id, want %d got %d", entry.ID, entry2.ID)
	}
	if entry.Amount != entry2.Amount {
		t.Errorf("failed deffrent amount id, want %d got %d", entry.Amount, entry2.Amount)
	}
}

func TestGetEntryByID(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	acc := createRandomAccount(user.Username)

	// test insert account
	newID, err := testRepo.InsertAccount(context.Background(), acc)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	acc.ID = newID

	// insert entry
	entry := models.Entry{
		AccountID: acc.ID,
		Amount:    acc.Balance,
	}

	newID, err = testRepo.InsertEntry(context.Background(), entry)
	if err != nil {
		t.Errorf("failed insert entry error:%s", err)
	}

	entry.ID = newID

	// get back entry
	entry2, err := testRepo.GetEntryByID(context.Background(), entry.ID)
	if err != nil {
		t.Errorf("failed get entry error:%s", err)
	}

	if entry.ID != entry2.ID {
		t.Errorf("failed deffrent entry id, want %d got %d", entry.ID, entry2.ID)
	}
	if entry.Amount != entry2.Amount {
		t.Errorf("failed deffrent amount id, want %d got %d", entry.Amount, entry2.Amount)
	}
}

func TestGetListEntries(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	acc := createRandomAccount(user.Username)
	// test insert account
	newID, err := testRepo.InsertAccount(context.Background(), acc)
	if err != nil {
		t.Errorf("failed insert account error:%s", err)
	}
	if newID < 1 {
		t.Errorf("failed renturn less than 1 newID:%d", newID)
	}

	acc.ID = newID

	for i := 0; i < 10; i++ {
		// update account
		acc.Balance = util.RandomBalance()

		err = testRepo.UpdateAccount(context.Background(), acc)
		if err != nil {
			t.Errorf("failed update account balance error:%s", err)
		}

		// insert entry
		entry := models.Entry{
			AccountID: acc.ID,
			Amount:    acc.Balance,
		}

		_, err = testRepo.InsertEntry(context.Background(), entry)
		if err != nil {
			t.Errorf("failed insert entry error:%s", err)
		}
	}

	listEntries, err := testRepo.GetListEntries(context.Background(), acc.ID, 10, 0)
	if err != nil {
		t.Errorf("failed get list entries error:%s", err)
	}

	if len(listEntries) < 10 {
		t.Errorf("failed len entries not 10 len:%d", len(listEntries))
	}
}

func TestInsertTransfer(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user1 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user1)
	if err != nil {
		t.Errorf("failed insert users 1 error:%s", err)
	}

	acc1 := createRandomAccount(user1.Username)

	user2 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err = testRepo.InsertUsers(context.Background(), user2)
	if err != nil {
		t.Errorf("failed insert users 2 error:%s", err)
	}

	acc2 := createRandomAccount(user2.Username)

	// insert 2 new account
	for _, acc := range []*models.Account{&acc1, &acc2} {
		newID, err := testRepo.InsertAccount(context.Background(), *acc)
		if err != nil {
			t.Errorf("failed insert account error:%s", err)
		}
		if newID < 1 {
			t.Errorf("failed renturn less than 1 newID:%d", newID)
		}

		acc.ID = newID
	}

	tf := models.Transfer{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomBalance(),
	}

	newID, err := testRepo.InsertTransfer(context.Background(), tf)
	if err != nil {
		t.Errorf("failed insert transfer error:%s", err)
	}

	tf.ID = newID

	tf2, err := testRepo.GetTransferByID(context.Background(), tf.ID)
	if err != nil {
		t.Errorf("failed get transfer error:%s", err)
	}

	if tf.ID != tf2.ID {
		t.Errorf("failed deffrent transfer_id want %d got %d", tf.ID, tf2.ID)
	}
	if tf.FromAccountID != tf2.FromAccountID {
		t.Errorf("failed deffrent from_account_id want %d got %d", tf.FromAccountID, tf2.FromAccountID)
	}
	if tf.ToAccountID != tf2.ToAccountID {
		t.Errorf("failed deffrent to_account_id want %d got %d", tf.ToAccountID, tf2.ToAccountID)
	}
	if tf.Amount != tf2.Amount {
		t.Errorf("failed deffrent amount want %d got %d", tf.Amount, tf2.Amount)
	}
}

func TestGetTransferByID(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user1 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user1)
	if err != nil {
		t.Errorf("failed insert users 1 error:%s", err)
	}

	acc1 := createRandomAccount(user1.Username)

	user2 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err = testRepo.InsertUsers(context.Background(), user2)
	if err != nil {
		t.Errorf("failed insert users 2 error:%s", err)
	}

	acc2 := createRandomAccount(user2.Username)

	// insert 2 new account
	for _, acc := range []*models.Account{&acc1, &acc2} {
		newID, err := testRepo.InsertAccount(context.Background(), *acc)
		if err != nil {
			t.Errorf("failed insert account error:%s", err)
		}
		if newID < 1 {
			t.Errorf("failed renturn less than 1 newID:%d", newID)
		}

		acc.ID = newID
	}

	tf := models.Transfer{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomBalance(),
	}

	newID, err := testRepo.InsertTransfer(context.Background(), tf)
	if err != nil {
		t.Errorf("failed insert transfer error:%s", err)
	}

	tf.ID = newID

	tf2, err := testRepo.GetTransferByID(context.Background(), tf.ID)
	if err != nil {
		t.Errorf("failed get transfer error:%s", err)
	}

	if tf.ID != tf2.ID {
		t.Errorf("failed deffrent transfer_id want %d got %d", tf.ID, tf2.ID)
	}
	if tf.FromAccountID != tf2.FromAccountID {
		t.Errorf("failed deffrent from_account_id want %d got %d", tf.FromAccountID, tf2.FromAccountID)
	}
	if tf.ToAccountID != tf2.ToAccountID {
		t.Errorf("failed deffrent to_account_id want %d got %d", tf.ToAccountID, tf2.ToAccountID)
	}
	if tf.Amount != tf2.Amount {
		t.Errorf("failed deffrent amount want %d got %d", tf.Amount, tf2.Amount)
	}
}

func TestGetListTransfers(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	user1 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), user1)
	if err != nil {
		t.Errorf("failed insert users 1 error:%s", err)
	}

	acc1 := createRandomAccount(user1.Username)

	user2 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err = testRepo.InsertUsers(context.Background(), user2)
	if err != nil {
		t.Errorf("failed insert users 2 error:%s", err)
	}

	acc2 := createRandomAccount(user2.Username)

	// insert 2 new account
	for _, acc := range []*models.Account{&acc1, &acc2} {
		newID, err := testRepo.InsertAccount(context.Background(), *acc)
		if err != nil {
			t.Errorf("failed insert account error:%s", err)
		}
		if newID < 1 {
			t.Errorf("failed renturn less than 1 newID:%d", newID)
		}

		acc.ID = newID
	}

	for i := 0; i < 10; i++ {
		tf := models.Transfer{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        util.RandomBalance(),
		}

		_, err := testRepo.InsertTransfer(context.Background(), tf)
		if err != nil {
			t.Errorf("failed insert transfer error:%s", err)
		}
	}

	listTf, err := testRepo.GetListTransfers(context.Background(), acc1.ID, acc2.ID, 10, 0)
	if err != nil {
		t.Errorf("failed to get list transfers error:%s", err)
	}

	if len(listTf) < 10 {
		t.Errorf("fialed len list transefer want %d got %d", 10, len(listTf))
	}
}

func TestInsertUsers(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	dataTest := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetUsersByUsername(context.Background(), dataTest.Username)
	if err != nil {
		t.Errorf("failed get users error:%s", err)
	}

	if dataTest.Username != data.Username {
		t.Errorf("failed deffrent username, wan't %s got %s", dataTest.Username, data.Username)
	}
	if dataTest.HashedPassword != data.HashedPassword {
		t.Errorf("failed deffrent hashed_password, wan't %s got %s", dataTest.HashedPassword, data.HashedPassword)
	}
	if dataTest.FullName != data.FullName {
		t.Errorf("failed deffrent full_name, wan't %s got %s", dataTest.FullName, data.FullName)
	}
	if dataTest.Email != data.Email {
		t.Errorf("failed deffrent email, wan't %s got %s", dataTest.Email, data.Email)
	}
}

func TestGetUsersByUsername(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	dataTest := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetUsersByUsername(context.Background(), dataTest.Username)
	if err != nil {
		t.Errorf("failed get users error:%s", err)
	}

	if dataTest.Username != data.Username {
		t.Errorf("failed deffrent username, wan't %s got %s", dataTest.Username, data.Username)
	}
	if dataTest.HashedPassword != data.HashedPassword {
		t.Errorf("failed deffrent hashed_password, wan't %s got %s", dataTest.HashedPassword, data.HashedPassword)
	}
	if dataTest.FullName != data.FullName {
		t.Errorf("failed deffrent full_name, wan't %s got %s", dataTest.FullName, data.FullName)
	}
	if dataTest.Email != data.Email {
		t.Errorf("failed deffrent email, wan't %s got %s", dataTest.Email, data.Email)
	}
}

func TestGetUsersByEmail(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	dataTest := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetUsersByEmail(context.Background(), dataTest.Email)
	if err != nil {
		t.Errorf("failed get users error:%s", err)
	}

	if dataTest.Username != data.Username {
		t.Errorf("failed deffrent username, wan't %s got %s", dataTest.Username, data.Username)
	}
	if dataTest.HashedPassword != data.HashedPassword {
		t.Errorf("failed deffrent hashed_password, wan't %s got %s", dataTest.HashedPassword, data.HashedPassword)
	}
	if dataTest.FullName != data.FullName {
		t.Errorf("failed deffrent full_name, wan't %s got %s", dataTest.FullName, data.FullName)
	}
	if dataTest.Email != data.Email {
		t.Errorf("failed deffrent email, wan't %s got %s", dataTest.Email, data.Email)
	}
}

func TestUpdateUsers(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	dataTest := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	// update
	dataTest.FullName = util.RandomOwner()

	err = testRepo.UpdateUsers(context.Background(), UpdateUserParam{
		Username: dataTest.Username,
		FullName: sql.NullString{
			String: dataTest.FullName,
			Valid:  true,
		},
	})
	if err != nil {
		t.Errorf("failed update users balance error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetUsersByUsername(context.Background(), dataTest.Username)
	if err != nil {
		t.Errorf("failed get account error:%s", err)
	}

	if dataTest.Username != data.Username {
		t.Errorf("failed deffrent username, wan't %s got %s", dataTest.Username, data.Username)
	}
	if dataTest.HashedPassword != data.HashedPassword {
		t.Errorf("failed deffrent hashed_password, wan't %s got %s", dataTest.HashedPassword, data.HashedPassword)
	}
	if dataTest.FullName != data.FullName {
		t.Errorf("failed deffrent full_name, wan't %s got %s", dataTest.FullName, data.FullName)
	}
	if dataTest.Email != data.Email {
		t.Errorf("failed deffrent email, wan't %s got %s", dataTest.Email, data.Email)
	}
}

func TestDeleteUsers(t *testing.T) {
	hashedPassword, _ := util.HashedPassword(util.RandomString(6))
	dataTest := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// test insert
	err := testRepo.InsertUsers(context.Background(), dataTest)
	if err != nil {
		t.Errorf("failed insert users error:%s", err)
	}

	// delete
	err = testRepo.DeleteUsers(context.Background(), dataTest.Username)
	if err != nil {
		t.Errorf("failed delete users error:%s", err)
	}

	// get back the data
	data, err := testRepo.GetUsersByUsername(context.Background(), dataTest.Username)
	if err != nil {
		t.Errorf("failed get users error:%s", err)
	}

	if data.Username != "" {
		t.Errorf("failed username must be empty")
	}
}

func TestGetListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		hashedPassword, _ := util.HashedPassword(util.RandomString(6))
		dataTest := models.Users{
			Username:       util.RandomOwner(),
			HashedPassword: hashedPassword,
			FullName:       util.RandomOwner(),
			Email:          util.RandomEmail(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// test insert
		err := testRepo.InsertUsers(context.Background(), dataTest)
		if err != nil {
			t.Errorf("failed insert users error:%s", err)
		}
	}

	listData, err := testRepo.GetListUsers(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("failed get list users error:%s", err)
	}

	if len(listData) < 10 {
		t.Errorf("failed len data not 10 len:%d", len(listData))
	}
}
