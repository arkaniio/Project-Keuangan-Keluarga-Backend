package model

import (
	"time"

	"github.com/google/uuid"
)

type Keuangan struct {
	Id                uuid.UUID `db:"id"`
	UserId            uuid.UUID `db:"user_id"`
	JenisTransaksi    string    `db:"jenis_transaksi"`
	JumlahPengeluaran int64     `db:"jumlah_pengeluaran"`
	JumlahPemasukan   int64     `db:"jumlah_pemasukan"`
	Kategori          string    `db:"kategori"`
	Tanggal           string    `db:"tanggal"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type PayloadKeuangan struct {
	Id                uuid.UUID `json:"id"`
	UserId            uuid.UUID `json:"user_id" validate:"required"`
	JenisTransaksi    string    `json:"jenis_transaksi" validate:"required"`
	JumlahPengeluaran int64     `json:"jumlah_pengeluaran" validate:"required"`
	JumlahPemasukan   int64     `json:"jumlah_pemasukan" validate:"required"`
	Kategori          string    `json:"kategori" validate:"required"`
	Tanggal           string    `json:"tanggal" validate:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PaylodUpdateKeuangan struct {
	Id                uuid.UUID `json:"id"`
	UserId            uuid.UUID `json:"user_id" validate:"required"`
	JenisTransaksi    *string   `json:"jenis_transaksi"`
	JumlahPengeluaran *int64    `json:"jumlah_pengeluaran"`
	JumlahPemasukan   *int64    `json:"jumlah_pemasukan"`
	Kategori          *string   `json:"kategori"`
	Tanggal           *string   `json:"tanggal"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StructureKeuanganWithUser struct {
	Id                uuid.UUID `json:"id"`
	UserId            uuid.UUID `json:"user_id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Profile_img       string    `json:"profile_img"`
	JenisTransaksi    string    `json:"jenis_transaksi"`
	JumlahPengeluaran int64     `json:"jumlah_pengeluaran"`
	JumlahPemasukan   int64     `json:"jumlah_pemasukan"`
	Kategori          string    `json:"kategori"`
	Tanggal           string    `json:"tanggal"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type KeuanganDataWithUser struct {
	Id                uuid.UUID `json:"id"`
	UserId            uuid.UUID `json:"user_id"`
	User              User      `json:"user"`
	JenisTransaksi    string    `json:"jenis_transaksi"`
	JumlahPengeluaran int64     `json:"jumlah_pengeluaran"`
	JumlahPemasukan   int64     `json:"jumlah_pemasukan"`
	Kategori          string    `json:"kategori"`
	Tanggal           string    `json:"tanggal"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
