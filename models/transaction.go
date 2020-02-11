package models

import (
	"context"
	"database/sql"
	"mini-wallet/libraries/api"
	"time"

	"github.com/google/uuid"
)

// Transaction : struct of Transaction
type Transaction struct {
	ID          string
	WalletID    string
	ReferenceID string
	CreatedBy   string
	Status      string
	Type        string
	Amount      float64
	Created     time.Time
}

// GetByReference : Get Transaction By Reference
func (u *Transaction) GetByReference(ctx context.Context, db *sql.DB) error {
	var isExist bool
	const q string = `SELECT true FROM transactions WHERE type = $1 AND reference_id = $2 `
	return db.QueryRowContext(ctx, q, u.Type, u.ReferenceID).Scan(&isExist)
}

// Save new transaction
func (u *Transaction) Save(ctx context.Context, tx *sql.Tx) error {
	u.ID = uuid.New().String()
	u.Created = time.Now().UTC()
	u.CreatedBy = ctx.Value(api.Ctx("auth")).(Customer).ID

	const query = `
		INSERT INTO transactions (transaction_id, wallet_id, reference_id, status, type, amount, created_by, created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, u.ID, u.WalletID, u.ReferenceID, u.Status, u.Type, u.Amount, u.CreatedBy, u.Created)
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return u.updateWallet(ctx, tx)
}

func (u *Transaction) updateWallet(ctx context.Context, tx *sql.Tx) error {
	var query string

	if u.Type == "in" {
		query = `UPDATE wallets SET balance = (balance + $1) WHERE wallet_id = $2`
	} else {
		query = `UPDATE wallets SET balance = (balance - $1) WHERE wallet_id = $2`
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Amount, u.WalletID)
	return err
}
