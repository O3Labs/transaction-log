package models

type TransactionLog struct {
	ID     string
	Source JSONB
	TXType string
	TXID   string
	Data   JSONB
}
