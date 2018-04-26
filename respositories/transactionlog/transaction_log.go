package transactionlog

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/o3labs/transaction-log/models"
)

type DatabaseRepository struct {
	DB *sql.DB
}

func (d *DatabaseRepository) SaveTransaction(tx models.TransactionLog) (*models.TransactionLog, error) {
	stmt, err := d.DB.Prepare(`
		insert into 
		transaction_log(txtype, txid, data, source) 
		values($1, $2, $3, $4) 
		returning id;`)

	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(tx.TXType, tx.TXID, tx.Data, tx.Source).Scan(&tx.ID)

	if err != nil {
		//we ignore error if it's a constraint one
		if strings.Contains(err.Error(), "transaction_log_txtype_txid_key") {
			return &tx, nil
		}
		return nil, fmt.Errorf("error %v", err)
	}
	return &tx, nil
}
