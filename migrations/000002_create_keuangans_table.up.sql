CREATE TABLE keuangans (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    jenis_transaksi VARCHAR(50) NOT NULL,
    jumlah_pengeluaran BIGINT NOT NULL,
    jumlah_pemasukan BIGINT NOT NULL,
    kategori VARCHAR(50) NOT NULL,
    tanggal DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);