package response

import (
	"mini-wallet/models"
)

// DepositResponse : format json response for deposit
type DepositResponse struct {
	Deposit struct {
		ID          string  `json:"id"`
		CustomerID  string  `json:"deposited_by"`
		Status      string  `json:"status"`
		DepositedAt string  `json:"deposited_at"`
		Amount      float64 `json:"amount"`
		ReferenceID string  `json:"reference_id"`
	} `json:"deposit"`
}

// Transform Transaction models to Deposit Response
func (u *DepositResponse) Transform(m models.Transaction) {
	if m.Type == "in" {
		u.Deposit.ID = m.ID
		u.Deposit.CustomerID = m.CreatedBy
		u.Deposit.Status = m.Status
		u.Deposit.DepositedAt = m.Created.Format("2006-01-02 15:04:05")
		u.Deposit.Amount = m.Amount
		u.Deposit.ReferenceID = m.ReferenceID
	}
}
