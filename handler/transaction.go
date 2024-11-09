package handler

import (
	"net/http"
	"task-golang-batch2/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionInterface interface {
	NewTransaction(*gin.Context)
	TransactionList(*gin.Context)
}

type transactionImplement struct {
	db *gorm.DB
}

// NewTransaction adalah handler untuk transaksi baru
func NewTransaction(db *gorm.DB) TransactionInterface {
	return &transactionImplement{
		db: db,
	}
}

// NewTransaction membuat record transaksi baru
func (t *transactionImplement) NewTransaction(c *gin.Context) {
	var payload model.Transaction

	// Bind JSON request ke payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validasi: cek apakah `account_id` disediakan dalam request atau context
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found in context"})
		return
	}

	payload.AccountID = accountID.(int64)

	// Set tanggal transaksi ke waktu saat ini jika tidak disediakan
	if payload.TransactionDate.IsZero() {
		payload.TransactionDate = time.Now()
	}

	// Buat record transaksi
	if err := t.db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction: " + err.Error()})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data":    payload,
	})
}

// TransactionList mengembalikan daftar transaksi berdasarkan `account_id`
// TransactionList mengembalikan daftar transaksi berdasarkan `account_id` dan filter tanggal
func (t *transactionImplement) TransactionList(c *gin.Context) {
	var transactions []model.Transaction

	// Ambil `account_id` dari context
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found in context"})
		return
	}

	// Ambil parameter filter tanggal dari query string (jika ada)
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	var startDate, endDate *time.Time

	// Jika start_date diberikan, konversi string ke time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr) // Format tanggal yang diinginkan: "YYYY-MM-DD"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
		startDate = &parsedStartDate
	}

	// Jika end_date diberikan, konversi string ke time.Time
	if endDateStr != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateStr) // Format tanggal yang diinginkan: "YYYY-MM-DD"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
		endDate = &parsedEndDate
	}

	// Siapkan query untuk mengambil transaksi berdasarkan account_id dan filter tanggal
	query := t.db.Where("account_id = ?", accountID)

	// Jika start_date disediakan, tambahkan kondisi WHERE untuk filter tanggal mulai
	if startDate != nil {
		query = query.Where("transaction_date >= ?", startDate)
	}

	// Jika end_date disediakan, tambahkan kondisi WHERE untuk filter tanggal akhir
	if endDate != nil {
		query = query.Where("transaction_date <= ?", endDate)
	}

	// Ambil 10 transaksi terakhir berdasarkan account_id dan filter tanggal (jika ada)
	if err := query.Order("transaction_date DESC").
		Limit(10).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions: " + err.Error()})
		return
	}

	// Respon sukses dengan data transaksi
	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}
