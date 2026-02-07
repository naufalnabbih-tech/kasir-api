package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

// Bertugas sebagai penghubung antara handler dan repository
type TransactionService struct {
	repo *repositories.TransactionRepository
}

// NewTransactionService membuat instance baru dari TransactionService
func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}
