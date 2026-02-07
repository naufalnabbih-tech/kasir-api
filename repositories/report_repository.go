package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetTodayReport() (*models.ReportResponse, error) {
	var report models.ReportResponse

	// Get total revenue dan total transaksi hari ini
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get produk terlaris hari ini
	err = r.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &report, nil
}

func (r *ReportRepository) GetReportByDateRange(startDate, endDate string) (*models.ReportResponse, error) {
	var report models.ReportResponse

	// Get total revenue dan total transaksi dalam range
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get produk terlaris dalam range
	err = r.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, startDate, endDate).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &report, nil
}
