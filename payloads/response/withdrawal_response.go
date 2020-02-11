package response

import (
	"mini-wallet/models"
)

// WithdrawalResponse : format json response for withdrawal
type WithdrawalResponse struct {
	Withdrawal struct {
		ID          string  `json:"id"`
		CustomerID  string  `json:"withdrawn_by"`
		Status      string  `json:"status"`
		WithdrawnAt string  `json:"withdrawn_at"`
		Amount      float64 `json:"balance"`
		ReferenceID string  `json:"reference_id"`
	} `json:"withdrawal"`
}

// Transform Transaction models to Withdrawal Response
func (u *WithdrawalResponse) Transform(m models.Transaction) {
	if m.Type == "out" {
		u.Withdrawal.ID = m.ID
		u.Withdrawal.CustomerID = m.CreatedBy
		u.Withdrawal.Status = m.Status
		u.Withdrawal.WithdrawnAt = m.Created.Format("2006-01-02 15:04:05")
		u.Withdrawal.Amount = m.Amount
		u.Withdrawal.ReferenceID = m.ReferenceID
	}
}
