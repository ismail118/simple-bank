package repository

import (
	"context"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/util"
	"log"
	"testing"
	"time"
)

func TestTransferTx(t *testing.T) {
	user1 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
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
		HashedPassword: "secret",
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

	log.Println(">> before:", acc1.Balance, acc2.Balance)

	newID, err := testStore.InsertAccount(context.Background(), acc1)
	if err != nil {
		t.Errorf("failed insert account 1 error:%s", err)
	}
	acc1.ID = newID

	newID, err = testStore.InsertAccount(context.Background(), acc2)
	if err != nil {
		t.Errorf("failed insert account 1 error:%s", err)
	}
	acc2.ID = newID

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	chErr := make(chan error)
	chRes := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			tf := models.Transfer{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			}
			result, err := testStore.TransferTx(context.Background(), tf)

			chErr <- err
			chRes <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-chErr
		res := <-chRes
		if err != nil {
			t.Fatalf("failed transfer error:%s", err)
		}

		// check result Transfer
		if res.Transfer.ID < 1 {
			t.Fatalf("failed transfer id is 0")
		}
		if res.Transfer.FromAccountID != acc1.ID {
			t.Fatalf("failed transfer from account id wan't %d got %d", acc1.ID, res.Transfer.FromAccountID)
		}
		if res.Transfer.ToAccountID != acc2.ID {
			t.Fatalf("failed transfer to account id wan't %d got %d", acc2.ID, res.Transfer.ToAccountID)
		}
		if res.Transfer.Amount != amount {
			t.Fatalf("faield transfer amount want %d got %d", amount, res.Transfer.Amount)
		}

		tf, err := testStore.GetTransferByID(context.Background(), res.Transfer.ID)
		if err != nil {
			t.Fatalf("failed get transfer error:%s", err)
		}
		if tf.ID < 1 {
			t.Fatalf("failed transfer with id: %d not found", res.Transfer.ID)
		}

		// check result entries
		if res.FromEntry.ID < 1 {
			t.Fatalf("failed from Entry id is 0")
		}
		if res.FromEntry.AccountID != acc1.ID {
			t.Fatalf("failed from entry account id want %d got %d", acc1.ID, res.FromEntry.AccountID)
		}
		if res.FromEntry.Amount != -amount {
			t.Fatalf("failed from entry amount want %d got %d", -amount, res.FromEntry.Amount)
		}

		if res.ToEntry.ID < 1 {
			t.Fatalf("failed to Entry id is 0")
		}
		if res.ToEntry.AccountID != acc2.ID {
			t.Fatalf("failed from entry account id want %d got %d", acc2.ID, res.ToEntry.AccountID)
		}
		if res.ToEntry.Amount != amount {
			t.Fatalf("failed from entry amount want %d got %d", amount, res.ToEntry.Amount)
		}

		// check accounts
		fromAccount := res.FromAccount
		if fromAccount.ID < 1 {
			t.Fatalf("failed from account id is 0")
		}
		toAccount := res.ToAccount
		if toAccount.ID < 1 {
			t.Fatalf("failed to account id is 0")
		}
		log.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		// check accounts balance
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		if diff1 != diff2 {
			t.Fatalf("failed deffrent balance diff1:%d dff2:%d", diff1, diff2)
		}
		if !(diff1 > 0) {
			t.Fatalf("failed diff1 must > 0, diff1:%d", diff1)
		}
		if !(diff1%amount == 0) {
			t.Fatalf("failed diff1 mod amount must be 0")
		}
		k := int(diff1 / amount)
		if !(k >= 1 && k <= n) {
			t.Fatalf("failed k")
		}
		_, ok := existed[k]
		if ok {
			t.Fatalf("faield %d should not exited in map", k)
		}
		existed[k] = true
	}

	// check the final updated balances
	updatedAcc1, err := testRepo.GetAccountByID(context.Background(), acc1.ID)
	if err != nil {
		t.Errorf("failed get update account 1 error:%s", err)
	}
	updatedAcc2, err := testRepo.GetAccountByID(context.Background(), acc2.ID)
	if err != nil {
		t.Errorf("failed get update account 2 error:%s", err)
	}

	log.Println(">> after:", updatedAcc1.Balance, updatedAcc2.Balance)

	balance1 := acc1.Balance - int64(n)*amount
	if balance1 != updatedAcc1.Balance {
		t.Errorf("failed updated balance account 1 want %d got %d", balance1, updatedAcc1.Balance)
	}
	balance2 := acc2.Balance + int64(n)*amount
	if balance2 != updatedAcc2.Balance {
		t.Errorf("failed updated balance account 2 want %d got %d", balance2, updatedAcc2.Balance)
	}

}

func TestTransferTxDeadlock(t *testing.T) {
	user1 := models.Users{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
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
		HashedPassword: "secret",
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

	log.Println(">> before:", acc1.Balance, acc2.Balance)

	newID, err := testStore.InsertAccount(context.Background(), acc1)
	if err != nil {
		t.Errorf("failed insert account 1 error:%s", err)
	}
	acc1.ID = newID

	newID, err = testStore.InsertAccount(context.Background(), acc2)
	if err != nil {
		t.Errorf("failed insert account 1 error:%s", err)
	}
	acc2.ID = newID

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	chErr := make(chan error)
	chRes := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 1 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}
		go func() {
			tf := models.Transfer{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			}
			res, err := testStore.TransferTx(context.Background(), tf)

			chErr <- err
			chRes <- res
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-chErr
		if err != nil {
			t.Fatalf("failed transfer error:%s", err)
		}
		res := <-chRes
		log.Println("tx:", res.FromAccount.Balance, res.ToAccount.Balance)

	}

	// check the final updated balances
	updatedAcc1, err := testRepo.GetAccountByID(context.Background(), acc1.ID)
	if err != nil {
		t.Errorf("failed get update account 1 error:%s", err)
	}
	updatedAcc2, err := testRepo.GetAccountByID(context.Background(), acc2.ID)
	if err != nil {
		t.Errorf("failed get update account 2 error:%s", err)
	}

	log.Println(">> after:", updatedAcc1.Balance, updatedAcc2.Balance)

	if acc1.Balance != updatedAcc1.Balance {
		t.Errorf("failed updated balance account 1 want %d got %d", acc1.Balance, updatedAcc1.Balance)
	}

	if acc2.Balance != updatedAcc2.Balance {
		t.Errorf("failed updated balance account 2 want %d got %d", acc2.Balance, updatedAcc2.Balance)
	}

}
