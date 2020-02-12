package schema

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add customers",
		Script: `
CREATE TABLE customers (
	customer_id	UUID,
	name 		TEXT NOT NULL,
	email 		TEXT NOT NULL UNIQUE,
	username    TEXT NOT NULL UNIQUE,
	password    TEXT NOT NULL,
	created 	TIMESTAMP NOT NULL DEFAULT NOW(),
	updated 	TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (customer_id)
);`,
	},
	{
		Version:     2,
		Description: "Add wallets",
		Script: `
CREATE TABLE wallets (
	wallet_id   UUID,
	customer_id UUID NOT NULL,
	status 		TEXT NOT NULL,
	balance 	DOUBLE PRECISION NOT NULL,
	created 	TIMESTAMP NOT NULL DEFAULT NOW(),
	updated 	TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (wallet_id),
	FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE CASCADE
);`,
	},
	{
		Version:     3,
		Description: "Add transactions",
		Script: `
CREATE TABLE transactions (
	transaction_id   	UUID,
	wallet_id 			UUID NOT NULL,
	reference_id		UUID NOT NULL UNIQUE,
	status 				TEXT NOT NULL,
	type 				CHAR(3) NOT NULL,
	amount 				DOUBLE PRECISION NOT NULL,
	created_by			UUID NOT NULL,
	created 			TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (transaction_id),
	FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id) ON DELETE CASCADE,
	FOREIGN KEY (created_by) REFERENCES customers(customer_id) ON DELETE CASCADE
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
