package models

import (
	"context"
	"database/sql"
	"mini-wallet/libraries/api"
	"time"

	"github.com/google/uuid"
)

// Wallet : struct of Wallet
type Wallet struct {
	ID         string
	CustomerID string
	Status     string
	Balance    float64
	EnableAt   time.Time
}

// Init new wallet
func (u *Wallet) Init(ctx context.Context, db *sql.DB) error {
	u.ID = uuid.New().String()
	u.Balance = 0
	u.Status = "enabled"
	u.EnableAt = time.Now().UTC()

	const query = `INSERT INTO wallets (wallet_id, customer_id, status, balance, created, updated) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, u.CustomerID, u.Status, u.Balance, u.EnableAt, u.EnableAt)
	return err
}

// Get wallet by customer
func (u *Wallet) Get(ctx context.Context, db *sql.DB) error {
	const q string = `SELECT wallet_id, customer_id, status, balance, updated FROM wallets WHERE customer_id = $1 `
	err := db.QueryRowContext(ctx, q, ctx.Value(api.Ctx("auth")).(Customer).ID).Scan(&u.ID, &u.CustomerID, &u.Status, &u.Balance, &u.EnableAt)

	if err == sql.ErrNoRows {
		err = api.ErrNotFound(err, "")
	}

	return err
}

// Disabled wallet by customer
func (u *Wallet) Disabled(ctx context.Context, db *sql.DB) error {
	const query = `UPDATE wallets SET status = $1, updated = $2 WHERE customer_id = $3`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, "disabled", time.Now().UTC(), ctx.Value(api.Ctx("auth")).(Customer).ID)
	if err != nil {
		return err
	}

	return u.Get(ctx, db)

}

// Enabled wallet by customer
func (u *Wallet) Enabled(ctx context.Context, db *sql.DB) error {
	const query = `UPDATE wallets SET status = $1, updated = $2 WHERE customer_id = $3`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, "enabled", time.Now().UTC(), ctx.Value(api.Ctx("auth")).(Customer).ID)
	if err != nil {
		return err
	}

	return u.Get(ctx, db)

}
