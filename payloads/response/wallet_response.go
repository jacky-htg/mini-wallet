package response

import (
	"mini-wallet/models"
)

// WalletResponse : format json response for wallet
type WalletResponse struct {
	Wallet struct {
		ID         string  `json:"id"`
		CustomerID string  `json:"owned_by"`
		Status     string  `json:"status"`
		EnableAt   string  `json:"enable_at"`
		Balance    float64 `json:"balance"`
	} `json:"wallet"`
}

// Transform Wallet models to Wallet Response
func (u *WalletResponse) Transform(m models.Wallet) {
	u.Wallet.ID = m.ID
	u.Wallet.CustomerID = m.CustomerID
	u.Wallet.Status = m.Status
	u.Wallet.EnableAt = m.EnableAt.Format("2006-01-02 15:04:05")
	u.Wallet.Balance = m.Balance
}
