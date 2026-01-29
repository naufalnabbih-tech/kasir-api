package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

// CategoryRepository adalah struct yang mengelola operasi database untuk tabel categories
// Kenapa menggunakan Repository pattern? Untuk memisahkan logika database dari business logic (separation of concerns)
type CategoryRepository struct {
	// db adalah pointer ke koneksi database yang akan digunakan untuk semua operasi
	// Kenapa menggunakan pointer (*sql.DB)? Agar tidak membuat copy koneksi database (lebih efisien dan hemat memori)
	db *sql.DB
}

// NewCategoryRepository adalah constructor function untuk membuat instance CategoryRepository
// Kenapa perlu constructor? Untuk dependency injection dan memastikan instance dibuat dengan benar
// Parameter db *sql.DB: menerima pointer koneksi database dari luar (dependency injection pattern)
// Kenapa return *CategoryRepository? Mengembalikan pointer agar lebih efisien (tidak copy struct)
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	// Membuat instance baru dengan field db diisi dari parameter
	// Kenapa menggunakan &CategoryRepository? Untuk mengembalikan alamat memori (pointer) bukan copy struct
	return &CategoryRepository{db: db}
}

// GetAll mengambil semua data kategori dari tabel categories
// Kenapa return ([]models.Category, error)? Pattern standar Go untuk mengembalikan data dan error
// Mengembalikan slice kategori dan error jika ada
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	// Query SQL untuk mengambil semua kategori dari tabel categories
	query := "SELECT id, name, description FROM categories"

	// Eksekusi query ke database dan simpan hasilnya dalam rows
	// Kenapa menggunakan Query()? Karena kita expect multiple rows (banyak kategori)
	rows, err := repo.db.Query(query)
	// Tangani error jika query gagal dieksekusi
	if err != nil {
		// Kembalikan nil dan error jika terjadi kesalahan
		// Kenapa return nil, err? Pattern Go untuk menandakan operasi gagal
		return nil, err
	}
	// Pastikan rows ditutup setelah selesai digunakan untuk menghindari memory leak
	// Kenapa defer? Agar Close() dipanggil otomatis sebelum function return (cleanup resource)
	defer rows.Close()

	// Buat slice kosong untuk menyimpan hasil kategori yang akan dikembalikan
	// Kenapa make([]models.Category, 0)? Untuk inisialisasi slice dengan kapasitas awal 0
	// Kenapa tidak nil? Agar JSON marshal menghasilkan [] bukan null
	categories := make([]models.Category, 0)
	// Loop melalui setiap baris hasil query
	// Kenapa rows.Next()? Untuk iterasi ke baris berikutnya, return false jika sudah habis
	for rows.Next() {
		// Deklarasi variabel temporary untuk menyimpan data kategori per baris
		// Kenapa var c models.Category? Untuk menyimpan hasil scan setiap iterasi
		var c models.Category
		// Scan data dari baris saat ini ke dalam struct Category
		// Kenapa pakai &c.ID, &c.Name, &c.Description? Scan butuh pointer untuk mengisi nilai
		// Kenapa urutan harus sama? Harus sesuai urutan kolom di SELECT query
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		// Cek apakah ada error saat scanning data
		if err != nil {
			// Kembalikan nil dan error jika scanning gagal
			// Kenapa langsung return? Tidak ada gunanya melanjutkan jika data corrupt
			return nil, err
		}
		// Tambahkan kategori yang sudah di-scan ke dalam slice categories
		// Kenapa append(categories, c)? Untuk menambahkan element ke slice secara dynamic
		categories = append(categories, c)
	}
	// Kembalikan slice categories yang sudah berisi semua data dan nil untuk error
	// Kenapa return categories, nil? nil menandakan tidak ada error (success)
	return categories, nil
}

// GetByID mengambil satu kategori berdasarkan ID
// Kenapa parameter id int? ID di database bertipe integer
// Kenapa return *models.Category? Pointer untuk menandakan bisa nil (not found) dan lebih efisien
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	// Query SQL untuk mengambil satu kategori berdasarkan ID dengan placeholder $1
	// Kenapa $1? Placeholder untuk prepared statement (mencegah SQL injection)
	// Kenapa WHERE id = $1? Filter untuk mengambil kategori dengan ID tertentu
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	// Deklarasi variabel untuk menyimpan hasil kategori yang akan di-scan
	var c models.Category
	// Eksekusi query dengan QueryRow (mengembalikan max 1 baris) dan langsung scan hasilnya
	// Kenapa QueryRow bukan Query? Karena kita expect maksimal 1 row berdasarkan ID (primary key)
	// Kenapa langsung .Scan()? QueryRow mengembalikan *Row yang bisa langsung di-scan
	// Kenapa parameter id? Nilai yang akan menggantikan placeholder $1
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	// Cek apakah data tidak ditemukan (ErrNoRows)
	// Kenapa cek sql.ErrNoRows khusus? Untuk membedakan "data tidak ada" vs "error database"
	if err == sql.ErrNoRows {
		// Kembalikan nil dan error custom jika kategori tidak ada di database
		// Kenapa return nil? Karena tidak ada data yang bisa dikembalikan
		return nil, errors.New("category not found")
	}
	// Cek error lainnya seperti error koneksi atau scanning
	if err != nil {
		// Kembalikan nil dan error asli jika terjadi kesalahan lain
		return nil, err
	}
	// Kembalikan pointer (alamat memori) ke struct Category dan nil untuk error
	// Kenapa &c? Mengambil alamat memori dari variabel c (return pointer)
	return &c, nil
}

// Create menambahkan kategori baru ke database
// Kenapa parameter *models.Category? Pointer agar bisa update field ID setelah insert
// Kenapa return error? Hanya perlu tahu berhasil atau gagal
func (repo *CategoryRepository) Create(category *models.Category) error {
	// Query SQL untuk menyisipkan kategori baru ke dalam tabel categories
	// Kenapa tidak INSERT id? Karena id auto-increment/serial, database yang generate
	// Kenapa RETURNING id? Untuk mendapatkan ID yang baru saja di-generate oleh database
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	// Eksekusi query dengan QueryRow untuk mendapatkan ID yang di-generate
	// Kenapa QueryRow? Karena RETURNING id mengembalikan 1 row berisi ID baru
	// Kenapa Scan(&category.ID)? Untuk menyimpan ID yang di-return ke struct category
	// Kenapa &category.ID? Pointer ke field ID agar bisa dimodifikasi (update by reference)
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	// Kembalikan error (nil jika sukses, ada nilai jika gagal)
	return err
}

// Update memperbarui data kategori yang sudah ada
// Kenapa parameter *models.Category? Menerima struct berisi data baru untuk di-update
// Kenapa return error? Untuk mengetahui apakah update berhasil atau gagal
func (repo *CategoryRepository) Update(category *models.Category) error {
	// Query SQL untuk memperbarui data kategori berdasarkan ID
	// Kenapa SET name, description? Field yang akan di-update (tidak termasuk id karena primary key)
	// Kenapa WHERE id = $3? Untuk memastikan hanya update kategori dengan ID tertentu
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	// Eksekusi query dengan Exec karena UPDATE tidak mengembalikan data, hanya result metadata
	// Kenapa Exec bukan Query? UPDATE tidak mengembalikan rows data, hanya info berapa row affected
	// Kenapa urutan parameter category.Name, Description, ID? Harus sesuai placeholder $1, $2, $3
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	// Cek apakah ada error saat eksekusi query (error koneksi, syntax, constraint, dll)
	if err != nil {
		// Kembalikan error jika query gagal dieksekusi
		return err
	}
	// Ambil jumlah baris yang terpengaruh oleh UPDATE untuk validasi
	// Kenapa perlu RowsAffected? Untuk memastikan kategori dengan ID tersebut benar-benar ada
	// Tanpa ini, update ID yang tidak ada akan return sukses (misleading!)
	rows, err := result.RowsAffected()
	// Cek apakah ada error saat mengambil RowsAffected
	if err != nil {
		// Kembalikan error jika gagal mendapatkan jumlah baris yang terpengaruh
		return err
	}
	// Jika rows == 0 artinya tidak ada baris yang di-update (ID tidak ditemukan di database)
	// Kenapa cek rows == 0? Untuk membedakan "sukses update" vs "ID tidak ada"
	if rows == 0 {
		// Kembalikan error custom untuk memberitahu bahwa kategori tidak ada
		// Kenapa error custom? Agar client tahu penyebab spesifik: data tidak ditemukan
		return errors.New("category not found")
	}
	// Kembalikan nil jika update berhasil (minimal 1 baris terpengaruh)
	// Kenapa return nil? nil = no error = success
	return nil
}

// Delete menghapus kategori dari database berdasarkan ID
// Kenapa parameter id int? Hanya butuh ID untuk menghapus, tidak perlu struct lengkap
// Kenapa return error? Untuk mengetahui apakah delete berhasil atau gagal
func (repo *CategoryRepository) Delete(id int) error {
	// Query SQL untuk menghapus kategori berdasarkan ID
	// Kenapa WHERE id = $1? Agar hanya menghapus kategori dengan ID tertentu (tidak semua data!)
	query := "DELETE FROM categories WHERE id = $1"
	// Eksekusi query dengan Exec karena DELETE tidak mengembalikan data, hanya result metadata
	// Kenapa Exec? DELETE tidak return rows data, hanya info berapa row deleted
	// Kenapa parameter id? Nilai yang akan menggantikan placeholder $1
	result, err := repo.db.Exec(query, id)
	// Cek apakah ada error saat eksekusi query (error koneksi, syntax, constraint, dll)
	// Contoh error: foreign key constraint (kategori masih dipakai di tabel lain)
	if err != nil {
		// Kembalikan error asli dari database
		return err
	}
	// Ambil jumlah baris yang terpengaruh oleh DELETE untuk validasi
	// Kenapa perlu RowsAffected? Untuk memastikan kategori dengan ID tersebut benar-benar ada
	// Tanpa ini, delete ID yang tidak ada akan return sukses (misleading!)
	rows, err := result.RowsAffected()
	// Cek apakah ada error saat mengambil RowsAffected
	if err != nil {
		// Kembalikan error jika gagal mendapatkan jumlah baris yang terpengaruh
		return err
	}
	// Jika rows == 0 artinya tidak ada baris yang dihapus (ID tidak ditemukan di database)
	// Kenapa cek rows == 0? Untuk membedakan "sukses delete" vs "ID tidak ada"
	if rows == 0 {
		// Kembalikan error custom untuk memberitahu bahwa kategori tidak ada
		// Kenapa error custom? Agar client tahu penyebab spesifik: data tidak ditemukan
		return errors.New("category not found")
	}
	// Kembalikan nil jika delete berhasil (minimal 1 baris terhapus)
	// Kenapa return nil? nil = no error = success
	return nil
}
