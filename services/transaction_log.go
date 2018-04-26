package services

import (
	"github.com/o3labs/transaction-log/database"
	"github.com/o3labs/transaction-log/models"
	"github.com/o3labs/transaction-log/respositories/transactionlog"
)

type TransactionLogInterface interface {
	SaveTransaction(tx models.TransactionLog) (*models.TransactionLog, error)
}
type TransactionLogService struct {
	TransactionLogInterface
}

func (t *TransactionLogService) SaveTransaction(tx models.TransactionLog) (*models.TransactionLog, error) {
	return t.TransactionLogInterface.SaveTransaction(tx)
}

func NewTransactionService() *TransactionLogService {
	db := database.Connect()
	r := &transactionlog.DatabaseRepository{DB: db}
	return &TransactionLogService{TransactionLogInterface: r}
}
