package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// ProductService menangani business logic untuk produk
// Bertugas sebagai penghubung antara handler dan repository
type ProductService struct {
	repo *repositories.ProductRepository
}

// NewProductService membuat instance baru dari ProductService
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll memanggil repository untuk mengambil semua produk
// Bisa ditambahkan validasi atau business logic di sini jika diperlukan
func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

// Create memvalidasi dan menyimpan produk baru melalui repository
// Di sini bisa ditambahkan validasi business logic seperti cek nama duplikat, validasi harga, dll
func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

// GetByID memanggil repository untuk mengambil produk berdasarkan ID
// Bisa ditambahkan business logic tambahan jika diperlukan
func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// Update memvalidasi dan memperbarui data produk melalui repository
// Bisa ditambahkan validasi seperti cek apakah produk ada, validasi perubahan data, dll
func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

// Delete menghapus produk melalui repository
// Bisa ditambahkan validasi seperti cek apakah produk sedang digunakan dalam transaksi, dll
func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
