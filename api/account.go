package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/models"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	acc := models.Account{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	newID, err := s.repo.InsertAccount(ctx, acc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	acc.ID = newID

	ctx.JSON(http.StatusAccepted, acc)
}

type getByIdRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getByIdRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	acc, err := s.repo.GetAccountByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if acc.ID < 1 {
		ctx.JSON(http.StatusNotFound, "account not found")
		return
	}

	ctx.JSON(http.StatusAccepted, acc)
}

type listRequest struct {
	Page int `json:"page" form:"page" binding:"required,min=1"`
	Size int `json:"size" form:"size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var req listRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accounts, err := s.repo.GetListAccounts(ctx,
		req.Size,
		(req.Page-1)*req.Size,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, accounts)
}

type updateAccountRequest struct {
	ID       int64  `json:"id" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (s *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.repo.GetAccountByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.ID < 1 {
		ctx.JSON(http.StatusNotFound, "account not found")
		return
	}

	account.Owner = req.Owner
	account.Balance = req.Balance
	account.Currency = req.Currency

	err = s.repo.UpdateAccount(ctx, account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, account)
}

func (s *Server) deleteAccount(ctx *gin.Context) {
	var req getByIdRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.repo.GetAccountByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.ID < 1 {
		ctx.JSON(http.StatusNotFound, "account not found")
		return
	}

	err = s.repo.DeleteAccount(ctx, account.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, "success delete account")
}

func (s *Server) getEntry(ctx *gin.Context) {
	var req getByIdRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := s.repo.GetEntryByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if entry.ID < 1 {
		ctx.JSON(http.StatusNotFound, "entry not found")
		return
	}

	ctx.JSON(http.StatusAccepted, entry)
}

type listEntryRequest struct {
	AccountID int64 `json:"account_id" form:"account_id" binding:"required,min=1"`
	Page      int   `json:"page" form:"page" binding:"required,min=1"`
	Size      int   `json:"size" form:"size" binding:"required,min=5,max=10"`
}

func (s *Server) listEntries(ctx *gin.Context) {
	var req listEntryRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	entries, err := s.repo.GetListEntries(ctx,
		req.AccountID,
		req.Size,
		(req.Page-1)*req.Size,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, entries)
}

func (s *Server) getTransfer(ctx *gin.Context) {
	var req getByIdRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := s.repo.GetTransferByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if transfer.ID < 1 {
		ctx.JSON(http.StatusNotFound, "transfer not found")
		return
	}

	ctx.JSON(http.StatusAccepted, transfer)
}

type listTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" form:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" form:"to_account_id" binding:"required,min=1"`
	Page          int   `json:"page" form:"page" binding:"required,min=1"`
	Size          int   `json:"size" form:"size" binding:"required,min=5,max=10"`
}

func (s *Server) listTransfer(ctx *gin.Context) {
	var req listTransferRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfers, err := s.repo.GetListTransfers(ctx,
		req.FromAccountID,
		req.ToAccountID,
		req.Size,
		(req.Page-1)*req.Size,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, transfers)
}

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,min=1"`
}

func (s *Server) transfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fAccount, err := s.repo.GetAccountByID(ctx, req.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if fAccount.ID < 1 {
		ctx.JSON(http.StatusNotFound, "from account not found")
		return
	}

	tAccount, err := s.repo.GetAccountByID(ctx, req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if tAccount.ID < 1 {
		ctx.JSON(http.StatusNotFound, "to account not found")
		return
	}

	tf := models.Transfer{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	res, err := s.store.TransferTx(ctx, tf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, res)
}
