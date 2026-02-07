package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	// prioritaskan env yang ada di os level
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Log config untuk debugging (jangan log password di production)
	fmt.Println("=== Configuration ===")
	fmt.Println("PORT:", config.Port)
	fmt.Println("DB_CONN exists:", config.DBConn != "")
	fmt.Println("=====================")

	// 1. Inisialisasi database terlebih dahulu
	fmt.Println("Attempting to connect to database...")
	fmt.Println("DB_CONN:", config.DBConn) // Log connection string (tanpa password)

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("ERROR: Failed to connect to database:", err)
		panic(err) // Panic agar Railway log error-nya
	}
	defer db.Close()
	fmt.Println("Database connected successfully!")

	// 2. Inisialisasi layer-layer aplikasi (Repository -> Service -> Handler)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// 3. Register routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/api/kategori", categoryHandler.HandleCategories)
	http.HandleFunc("/api/kategori/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	//  localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Test database connection
		err := db.Ping()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Database connection failed",
				"status":  "ERROR",
				"error":   err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "API Running",
			"status":   "OK",
			"database": "connected",
		})
	})

	// 4. Start server (ini harus paling akhir)
	addr := "0.0.0.0:" + config.Port
	fmt.Println("===========================================")
	fmt.Println("Server starting on", addr)
	fmt.Println("Health check: http://" + addr + "/health")
	fmt.Println("===========================================")

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("ERROR: Failed to start server:", err)
		panic(err)
	}
}
