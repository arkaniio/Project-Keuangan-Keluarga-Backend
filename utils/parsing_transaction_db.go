package utils

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

func AddTransaction(db *sqlx.DB, ctx context.Context) (*sqlx.Tx, error) {

	tx_options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	tx, err := db.BeginTxx(ctx, tx_options)
	if err != nil {
		return nil, errors.New("Failed to settings the transactions")
	}

	return tx, nil

}
