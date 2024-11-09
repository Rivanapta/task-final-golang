package model

import "time"

type Transaction struct {
	TransactionID         int64     `json:"transaction_id" gorm:"primaryKey;autoIncrement"`
	TransactionCategoryID *int64    `json:"transaction_category_id"` // Optional category for bonus
	AccountID             int64     `json:"account_id"`
	FromAccountID         *int64    `json:"from_account_id,omitempty"`
	ToAccountID           *int64    `json:"to_account_id,omitempty"`
	Amount                int64     `json:"amount"`
	TransactionDate       time.Time `json:"transaction_date" gorm:"autoCreateTime"`
}

// Menambahkan metode TableName untuk menyesuaikan nama tabel di database
func (Transaction) TableName() string {
	return "transaction" // nama tabel yang sesuai di database
}
