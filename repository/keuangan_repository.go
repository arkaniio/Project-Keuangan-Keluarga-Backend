package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-keuangan-keluarga/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type KeuanganRepository interface {
	CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error
	DeleteDataKeuangan(ctx context.Context, id uuid.UUID) error
}

type repoKeuangan struct {
	db *sqlx.DB
}

func NewKeuanganRepository(db *sqlx.DB) KeuanganRepository {
	return &repoKeuangan{db: db}
}

func (r *repoKeuangan) CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions")
	}
	defer tx.Rollback()

	query := `
		INSERT INTO keuangans(id, user_id, jenis_transaksi, jumlah_pengeluaran, jumlah_pemasukan, kategori, tanggal, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	rows, err := tx.ExecContext(ctx, query, keuangan.Id, keuangan.UserId, keuangan.JenisTransaksi, keuangan.JumlahPengeluaran, keuangan.JumlahPemasukan, keuangan.Kategori, keuangan.Tanggal, keuangan.CreatedAt, keuangan.UpdatedAt)
	if err != nil {
		return errors.New("Failed to exec the context" + err.Error())
	}

	rows_dected, _ := rows.RowsAffected()
	if rows_dected == 0 {
		return errors.New("Failed to detect the rows affected")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction")
	}

	return nil

}

func (r *repoKeuangan) DeleteDataKeuangan(ctx context.Context, Id uuid.UUID) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := r.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions")
	}
	defer tx.Rollback()

	query := `
		DELETE FROM keuangans WHERE id = $1;
	`

	rows, err := tx.ExecContext(ctx, query, Id)
	if err != nil {
		return errors.New("Failed to exec the context" + err.Error())
	}

	rows_dected, _ := rows.RowsAffected()
	if rows_dected == 0 {
		return errors.New("Failed to detect the rows affected")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction")
	}

	return nil

}
