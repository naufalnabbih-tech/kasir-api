package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository membuat instance baru dari TransactionRepository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin() // Menandakan memakai transaksi
	if err != nil {            // Jika error langsung return error
		return nil, err
	}
	defer tx.Rollback() // Jika ada error di tengah-tengah, maka rollback.

	//inisialisasi sub total -> jumlah total keseluruhan transaksi
	totalAmount := 0
	//inisialisasi modelling detail transaksi -> untuk insert ke db
	details := make([]models.TransactionDetails, 0)
	//loop setiap item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		//get product untuk mendapatkan harga
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d NOT FOUND", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
		//hitung current total = quantity * harga
		//ditambah ke dalam subtotal
		subtotal := price * item.Quantity
		totalAmount += subtotal
		//kurangi jumlah stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		//itemnya dimasukan ke transaction details
		details = append(details, models.TransactionDetails{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	//insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}
	//insert transaction details
	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil { //Jika semua proses berhasil, commit transaksi
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}
