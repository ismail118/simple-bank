package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
	"net/http"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	acc := models.Account{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	user, err := s.repo.GetUsersByUsername(ctx, acc.Owner)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Username == "" {
		ctx.JSON(http.StatusForbidden, fmt.Sprintf("user with username %s not exists", acc.Owner))
		return
	}

	a, err := s.repo.GetAccountByOwnerAndCurrency(ctx, acc.Owner, acc.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if a.ID > 1 {
		ctx.JSON(http.StatusForbidden, fmt.Sprintf("account with owner %s and %s alredy exists", acc.Owner, acc.Currency))
		return
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

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if acc.Owner != authPayload.Username {
		err = errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := s.repo.GetListAccounts(ctx,
		authPayload.Username,
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

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

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

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	account, err := s.repo.GetAccountByID(ctx, entry.AccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("entry doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.repo.GetAccountByID(ctx, req.AccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("entries doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	fAccount, err := s.repo.GetAccountByID(ctx, transfer.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tAccount, err := s.repo.GetAccountByID(ctx, transfer.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if (fAccount.Owner != authPayload.Username) && (tAccount.Owner != authPayload.Username) {
		err = errors.New("transfer doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

	fAccount, err := s.repo.GetAccountByID(ctx, req.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tAccount, err := s.repo.GetAccountByID(ctx, req.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if (fAccount.Owner != authPayload.Username) && (tAccount.Owner != authPayload.Username) {
		err = errors.New("transfer doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
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
	if fAccount.Currency != req.Currency {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("account %d mismatch: %s vs %s", fAccount.ID, fAccount.Currency, req.Currency))
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
	if fAccount.Currency != req.Currency {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("account %d mismatch: %s vs %s", tAccount.ID, tAccount.Currency, req.Currency))
		return
	}

	// authorization
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fAccount.Owner {
		err = errors.New("from account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user := models.Users{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	u, err := s.repo.GetUsersByUsername(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if u.Username != "" {
		ctx.JSON(http.StatusForbidden, fmt.Sprintf("user with username %s is exists", u.Username))
		return
	}

	u, err = s.repo.GetUsersByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if u.Username != "" {
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("email %s already being userd", req.Email))
		return
	}

	err = s.repo.InsertUsers(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	u, err = s.repo.GetUsersByUsername(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, u)
}

type getByUsernameRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (s *Server) getUsers(ctx *gin.Context) {
	var req getByUsernameRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Username == "" {
		ctx.JSON(http.StatusNotFound, "users not found")
		return
	}

	// authorization
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != user.Username {
		err = errors.New("user doesn't belong to user login")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, user)
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req listRequest

	err := ctx.BindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := s.repo.GetListUsers(ctx,
		req.Size,
		(req.Page-1)*req.Size,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, users)
}

type updateUsersRequest struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) updateUsers(ctx *gin.Context) {
	var req updateUsersRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Username == "" {
		ctx.JSON(http.StatusNotFound, "users not found")
		return
	}

	if user.Email != req.Email {
		u, err := s.repo.GetUsersByEmail(ctx, req.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if u.Username != "" {
			ctx.JSON(http.StatusNotFound, fmt.Sprintf("email %s already being userd", req.Email))
			return
		}
	}

	// authorization
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != user.Username {
		err = errors.New("user doesn't belong to user login")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	user.FullName = req.FullName
	user.Email = req.Email

	err = s.repo.UpdateUsers(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, req)
}

func (s *Server) deleteUsers(ctx *gin.Context) {
	var req getByUsernameRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Username == "" {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != user.Username {
		err = errors.New("user doesn't belong to user login")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	err = s.repo.DeleteUsers(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, "success delete users")
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        models.Users
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.repo.GetUsersByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Username == "" {
		ctx.JSON(http.StatusNotFound, "username not found")
		return
	}

	err = util.ComparePassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "wrong password")
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: accessToken,
		User:        user,
	}

	ctx.JSON(http.StatusAccepted, resp)
}
