package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project-keuangan-keluarga/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type KeuanganRepository interface {
	CreateNewKeuangan(ctx context.Context, keuangan *model.Keuangan) error
	DeleteDataKeuangan(ctx context.Context, id uuid.UUID) error
	UpdateKeuangan(ctx context.Context, id uuid.UUID, payload model.PaylodUpdateKeuangan) error
	GetAllKeuangans(ctx context.Context) ([]model.KeuanganDataWithUser, error)
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

func (r *repoKeuangan) UpdateKeuangan(ctx context.Context, id uuid.UUID, payload model.PaylodUpdateKeuangan) error {

	options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	tx, err := r.db.BeginTxx(ctx, options)
	if err != nil {
		return errors.New("Failed to setup the transactions")
	}
	defer tx.Rollback()

	var args []interface{}
	argsID := 1
	var settings []string

	if payload.JenisTransaksi != nil {
		settings = append(settings, fmt.Sprintf("jenis_transaksi=$%d", argsID))
		args = append(args, *payload.JenisTransaksi)
		argsID++
	}
	if payload.JumlahPengeluaran != nil {
		settings = append(settings, fmt.Sprintf("jumlah_pengeluaran=$%d", argsID))
		args = append(args, *payload.JumlahPengeluaran)
		argsID++
	}
	if payload.JumlahPemasukan != nil {
		settings = append(settings, fmt.Sprintf("jumlah_pemasukan=$%d", argsID))
		args = append(args, *payload.JumlahPemasukan)
		argsID++
	}
	if payload.Kategori != nil {
		settings = append(settings, fmt.Sprintf("kategori=$%d", argsID))
		args = append(args, *payload.Kategori)
		argsID++
	}
	if payload.Tanggal != nil {
		settings = append(settings, fmt.Sprintf("tanggal=$%d", argsID))
		args = append(args, *payload.Tanggal)
		argsID++
	}

	settings = append(settings, fmt.Sprintf("updated_at=$%d", argsID))
	args = append(args, time.Now().UTC())
	argsID++

	full_query := fmt.Sprintf("UPDATE keuangans SET %s WHERE id = $%d", strings.Join(settings, ", "), argsID)
	args = append(args, id)

	rows, err := tx.ExecContext(ctx, full_query, args...)
	if err != nil {
		return errors.New("Failed to exec the query!" + err.Error())
	}

	result, err := rows.RowsAffected()
	if err != nil {
		return errors.New("Failed to get the rows affected!")
	}

	if result == 0 {
		return errors.New("No one data expected!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to update the query and commit the query!")
	}

	return nil

}

func (r *repoKeuangan) GetAllKeuangans(ctx context.Context) ([]model.KeuanganDataWithUser, error) {

	query := `
		SELECT k.id, k.user_id, k.jenis_transaksi, k.jumlah_pengeluaran, k.jumlah_pemasukan, k.kategori, k.tanggal, k.created_at, k.updated_at
		FROM keuangans k
		JOIN users u ON k.user_id = u.id;
	`

	scan_query, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get the rows from the result!")
		}
		return nil, errors.New("Failed to detect the result of the query!")
	}
	defer scan_query.Close()

	var keuangan []model.KeuanganDataWithUser
	for scan_query.Next() {
		var keuangan_data model.StructureKeuanganWithUser
		if err := scan_query.StructScan(keuangan_data); err != nil {
			return nil, errors.New("Failed to get the struct of keuangan and user data!" + err.Error())
		}

		keuangans_data := model.KeuanganDataWithUser{
			Id:     keuangan_data.Id,
			UserId: keuangan_data.UserId,
			User: model.User{
				Id:          keuangan_data.UserId,
				Name:        keuangan_data.Name,
				Email:       keuangan_data.Email,
				Profile_img: keuangan_data.Profile_img,
			},
			JenisTransaksi:    keuangan_data.JenisTransaksi,
			JumlahPengeluaran: keuangan_data.JumlahPengeluaran,
			JumlahPemasukan:   keuangan_data.JumlahPemasukan,
			Kategori:          keuangan_data.Kategori,
			Tanggal:           keuangan_data.Tanggal,
			CreatedAt:         keuangan_data.CreatedAt,
			UpdatedAt:         keuangan_data.UpdatedAt,
		}

		keuangan = append(keuangan, keuangans_data)

	}

	return keuangan, nil

}
