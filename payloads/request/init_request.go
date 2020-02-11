package request

import (
	"mini-wallet/models"
)

// InitRequest : format json request for init wallet
type InitRequest struct {
	CustomerID string `json:"customer_id"`
}

// Transform to Wallet Model
func (u *InitRequest) Transform() models.Wallet {
	var m models.Wallet
	m.CustomerID = u.CustomerID

	return m
}
