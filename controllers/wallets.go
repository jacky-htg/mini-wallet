package controllers

import (
	"database/sql"
	"errors"
	"log"
	"mini-wallet/libraries/api"
	"mini-wallet/models"
	"mini-wallet/payloads/request"
	"mini-wallet/payloads/response"
	"net/http"
)

// Wallets : struct for set Wallets Dependency Injection
type Wallets struct {
	Db  *sql.DB
	Log *log.Logger
}

// Init : http handler for init wallet
func (u *Wallets) Init(w http.ResponseWriter, r *http.Request) {
	initRequest := new(request.InitRequest)
	err := api.Decode(r, &initRequest, false)
	if err != nil {
		u.Log.Printf("error decode init wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	if len(initRequest.CustomerID) <= 0 {
		err = api.ErrBadRequest(errors.New("customer_id is required"), "")
		u.Log.Printf("error validate init wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	wallet := initRequest.Transform()
	err = wallet.Init(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error init wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	res := new(response.WalletResponse)
	res.Transform(wallet)
	api.ResponseOK(w, res, http.StatusOK)
}

// View wallet
func (u *Wallets) View(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	err := wallet.Get(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	if wallet.Status == "disabled" {
		err = api.ErrForbidden(errors.New("your wallet is disabled"), "your wallet is disabled")
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	res := new(response.WalletResponse)
	res.Transform(wallet)
	api.ResponseOK(w, res, http.StatusOK)
}

// Enabled wallet
func (u *Wallets) Enabled(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	err := wallet.Enabled(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error enabled wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	res := new(response.WalletResponse)
	res.Transform(wallet)
	api.ResponseOK(w, res, http.StatusOK)
}

// Disabled wallet
func (u *Wallets) Disabled(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	err := wallet.Disabled(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error disabled wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	res := new(response.WalletResponse)
	res.Transform(wallet)
	api.ResponseOK(w, res, http.StatusOK)
}

// Deposit to wallet
func (u *Wallets) Deposit(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	err := wallet.Get(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	if wallet.Status == "disabled" {
		err = api.ErrForbidden(errors.New("your wallet is disabled"), "your wallet is disabled")
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	var transactionRequest request.TransactionRequest
	err = transactionRequest.Transform(r)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	tr := transactionRequest.TransformToModel()
	tr.Type = "in"
	tr.Status = "success"
	tr.WalletID = wallet.ID

	err = tr.GetByReference(r.Context(), u.Db)
	if err != sql.ErrNoRows {
		if err == nil {
			err = api.ErrBadRequest(errors.New("reference_id have been used"), "reference_id have been used")
		}

		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.BeginTx(r.Context(), &sql.TxOptions{})
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = tr.Save(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("error deposit wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	res := new(response.DepositResponse)
	res.Transform(tr)
	api.ResponseOK(w, res, http.StatusOK)
}

// Withdrawal from wallet
func (u *Wallets) Withdrawal(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	err := wallet.Get(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	if wallet.Status == "disabled" {
		err = api.ErrForbidden(errors.New("your wallet is disabled"), "your wallet is disabled")
		u.Log.Printf("error get wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	var transactionRequest request.TransactionRequest
	err = transactionRequest.Transform(r)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	if transactionRequest.Amount > wallet.Balance {
		err = api.ErrForbidden(errors.New("your balance is not enough"), "your balance is not enough")
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	tr := transactionRequest.TransformToModel()
	tr.Type = "out"
	tr.Status = "success"
	tr.WalletID = wallet.ID

	err = tr.GetByReference(r.Context(), u.Db)
	if err != sql.ErrNoRows {
		if err == nil {
			err = api.ErrBadRequest(errors.New("reference_id have been used"), "reference_id have been used")
		}

		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	tx, err := u.Db.BeginTx(r.Context(), &sql.TxOptions{})
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	err = tr.Save(r.Context(), tx)
	if err != nil {
		tx.Rollback()
		u.Log.Printf("error withdrawal wallet: %s", err)
		api.ResponseError(w, err)
		return
	}

	tx.Commit()

	res := new(response.WithdrawalResponse)
	res.Transform(tr)
	api.ResponseOK(w, res, http.StatusOK)
}
