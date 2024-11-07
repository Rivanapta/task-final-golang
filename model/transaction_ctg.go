package model

type TransactionCategories struct {
	TransactionCatID int64  `json:"transaction_category_id" gorm:"column:transaction_category_id;primaryKey;autoIncrement"`
	Name             string `json:"name" gorm:"column:name"`
}
