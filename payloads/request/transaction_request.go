package request

import (
	"errors"
	"mini-wallet/libraries/api"
	"mini-wallet/models"
	"net/http"
	"strconv"
)

// TransactionRequest : format json request for transaction
type TransactionRequest struct {
	Amount      float64
	ReferenceID string
}

// Transform form to TransactionRequest. Include Form Validate.
func (u *TransactionRequest) Transform(r *http.Request) error {
	if len(r.FormValue("reference_id")) <= 0 {
		return api.ErrBadRequest(errors.New("reference_id is required"), "reference_id is required")
	}

	if len(r.FormValue("amount")) <= 0 {
		return api.ErrBadRequest(errors.New("amount is required"), "amount is required")
	}

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		return err
	}

	if amount < 0 {
		return api.ErrBadRequest(errors.New("please suplay valid amount"), "please suplay valid amount")
	}

	u.Amount = float64(amount)
	u.ReferenceID = r.FormValue("reference_id")

	return nil
}

// TransformToModel : Transform TransactionRequest to Transaction Model
func (u *TransactionRequest) TransformToModel() models.Transaction {
	var m models.Transaction
	m.Amount = u.Amount
	m.ReferenceID = u.ReferenceID

	return m
}
