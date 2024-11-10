package handler

import (
	"errors"
	"net/http"
	"strconv"
	"task-golang-batch2/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountInterface interface {
	Create(*gin.Context)
	Read(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	List(*gin.Context)
	TopUp(*gin.Context)
	My(*gin.Context)
	Transfer(*gin.Context)
	Balance(*gin.Context)
	Mutation(*gin.Context)
}

type accountImplement struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) AccountInterface {
	return &accountImplement{
		db: db,
	}
}

func (a *accountImplement) Create(c *gin.Context) {
	payload := model.Account{}

	// bind JSON Request to payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Create data
	result := a.db.Create(&payload)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload,
	})
}

func (a *accountImplement) Read(c *gin.Context) {
	var account model.Account

	// get id from url account/read/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and put to account model
	if err := a.db.First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

func (a *accountImplement) Update(c *gin.Context) {
	payload := model.Account{}

	// bind JSON Request to payload
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// get id from url account/update/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and put to account model
	account := model.Account{}
	result := a.db.First(&account, "account_id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Update data
	account.Name = payload.Name
	account.Balance = payload.Balance
	a.db.Save(account)

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Update success",
	})
}

func (a *accountImplement) Delete(c *gin.Context) {
	// get id from url account/delete/5, 5 will be the id
	id := c.Param("id")

	// Find first data based on id and delete it
	if err := a.db.Where("account_id = ?", id).Delete(&model.Account{}).Error; err != nil {
		// No data found and deleted
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
		"data": map[string]string{
			"account_id": id,
		},
	})
}

func (a *accountImplement) List(c *gin.Context) {
	// Prepare empty result
	var accounts []model.Account

	// Find and get all accounts data and put to &accounts
	if err := a.db.Find(&accounts).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}

func (a *accountImplement) My(c *gin.Context) {
	var account model.Account
	// get account_id from middleware auth
	accountID := c.GetInt64("account_id")

	// Find first data based on account_id given
	if err := a.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

// Handler for "POST /account/topup"

func (a *accountImplement) TopUp(c *gin.Context) {
	// Ambil account_id dan amount dari form body
	accountIDStr := c.PostForm("account_id")
	amountStr := c.PostForm("amount")

	// Konversi account_id ke int64
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64) // Pastikan menggunakan ParseInt untuk int64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account_id"})
		return
	}

	// Konversi amount ke int64
	amount, err := strconv.ParseInt(amountStr, 10, 64) // Menggunakan ParseInt untuk amount
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Pastikan amount > 0
	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	// Ambil akun berdasarkan account_id (gunakan account_id sebagai kolom)
	var account model.Account
	err = a.db.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		// Jika account tidak ditemukan
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Mulai transaksi untuk top-up
	err = a.db.Transaction(func(tx *gorm.DB) error {
		// Update saldo akun dengan menambah jumlah top-up
		if err := tx.Model(&model.Account{}).Where("account_id = ?", accountID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// Buat entri transaksi untuk top-up
		transaction := model.Transaction{
			AccountID:             accountID,
			Amount:                amount,
			TransactionCategoryID: nil, // Bisa disesuaikan dengan kategori
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Top-Up Transaction Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance top-up completed successfully"})
}

func (a *accountImplement) Transfer(c *gin.Context) {
	// Ambil account_id dari middleware auth (harus ada casting ke int64)
	fromAccountIDInterface, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found"})
		return
	}
	fromAccountID := fromAccountIDInterface.(int64)

	// Ambil to_account_id dan amount dari form body
	toAccountIDStr := c.PostForm("to_account_id")
	amountStr := c.PostForm("amount")

	// Konversi toAccountID dan amount ke int64
	toAccountID, err := strconv.ParseInt(toAccountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_account_id"})
		return
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Mulai transaksi transfer
	err = a.db.Transaction(func(tx *gorm.DB) error {
		// Ambil akun pengirim
		var fromAccount model.Account
		if err := tx.Where("account_id = ?", fromAccountID).First(&fromAccount).Error; err != nil {
			return err
		}

		// Pastikan saldo mencukupi
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Kurangi saldo dari akun pengirim
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
			return err
		}

		// Tambah saldo ke akun penerima
		var toAccount model.Account
		if err := tx.Where("account_id = ?", toAccountID).First(&toAccount).Error; err != nil {
			return err
		}
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
			return err
		}

		// Buat catatan transaksi
		transaction := model.Transaction{
			AccountID:     fromAccountID,
			FromAccountID: &fromAccountID,
			ToAccountID:   &toAccountID,
			Amount:        amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transfer completed successfully"})
}

// Handler for "POST /account/balance"
func (h *accountImplement) Balance(c *gin.Context) {
	// Mendapatkan account_id dari context
	accountID := c.GetInt64("account_id")

	// Periksa jika account_id adalah 0, yang berarti tidak valid
	if accountID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not provided or invalid"})
		return
	}

	var account model.Account

	// Query untuk mengambil saldo berdasarkan account_id
	err := h.db.Select("balance").Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		// Jika akun tidak ditemukan, berikan error yang lebih spesifik
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

func (a *accountImplement) Mutation(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not authorized"})
		return
	}

	var transactions []model.Transaction
	var startDate, endDate time.Time
	var err error

	// Mengambil start_date dan end_date dari query parameter
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	// Jika ada start_date, parse menjadi time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
	}

	// Jika ada end_date, parse menjadi time.Time
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
	}

	// Ambil transaksi yang melibatkan akun sebagai pengirim atau penerima dan sesuaikan dengan tanggal
	query := a.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID)

	if !startDate.IsZero() {
		query = query.Where("transaction_date >= ?", startDate)
	}

	if !endDate.IsZero() {
		query = query.Where("transaction_date <= ?", endDate)
	}

	err = query.Order("transaction_date DESC").Limit(10).Find(&transactions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// func (a *accountImplement) Mutation(c *gin.Context) {
// 	accountID, exists := c.Get("account_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not authorized"})
// 		return
// 	}

// 	var transactions []model.Transaction

// 	// Ambil transaksi yang melibatkan akun sebagai pengirim atau penerima
// 	err := a.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
// 		Find(&transactions).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": transactions})
// }
