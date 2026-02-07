package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

// ProductRepository mengelola operasi database untuk tabel products
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository membuat instance baru dari ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll mengambil semua data produk dari tabel products
// Mengembalikan slice dari Product dan error jika ada
func (repo *ProductRepository) GetAll(nameFilter string) ([]models.Product, error) {
	query := `
	SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, '') as category_name
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	`
	args := []interface{}{}
	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetByID mengambil satu produk berdasarkan ID dari database
// Mengembalikan pointer ke Product dan error jika produk tidak ditemukan
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
	SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, '') as category_name
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	WHERE p.id = $1`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Create menambahkan produk baru ke database
// Mengisi field ID pada product dengan ID yang di-generate oleh database
func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// Update memperbarui data produk yang sudah ada di database
// Mengembalikan error jika produk dengan ID tersebut tidak ditemukan
func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Delete menghapus produk dari database berdasarkan ID
// Mengembalikan error jika produk dengan ID tersebut tidak ditemukan
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
