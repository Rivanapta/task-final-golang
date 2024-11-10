package main

import (
	"log"
	"os"
	"task-golang-batch2/handler"
	"task-golang-batch2/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Database
	db := NewDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	// secret-key
	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("SIGNING_KEY not set in environment")
	}

	r := gin.Default()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"}, // Ganti dengan domain yang diperbolehkan
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Apply CORS as a Gin middleware
	r.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	})

	// grouping route with /auth
	authHandler := handler.NewAuth(db, []byte(signingKey))
	authRoute := r.Group("/auth")
	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/upsert", authHandler.Upsert)
	// Tambahkan route baru untuk /change-password dengan menggunakan middleware AuthMiddleware
	authRoute.POST("/change-password", middleware.AuthMiddleware(signingKey), authHandler.ChangePassword)

	// grouping route with /account
	accountHandler := handler.NewAccount(db)
	accountRoutes := r.Group("/account")
	accountRoutes.POST("/create", accountHandler.Create)
	accountRoutes.GET("/read/:id", accountHandler.Read)
	accountRoutes.PATCH("/update/:id", accountHandler.Update)
	accountRoutes.DELETE("/delete/:id", accountHandler.Delete)
	accountRoutes.GET("/list", accountHandler.List)
	accountRoutes.POST("/topup", accountHandler.TopUp)
	accountRoutes.POST("/transfer", middleware.AuthMiddleware(signingKey), accountHandler.Transfer)
	accountRoutes.GET("/mutation", middleware.AuthMiddleware(signingKey), accountHandler.Mutation)
	accountRoutes.GET("/balance", middleware.AuthMiddleware(signingKey), accountHandler.Balance)

	accountRoutes.GET("/my", middleware.AuthMiddleware(signingKey), accountHandler.My)

	// grouping route with /transactionCategories
	transacttionCTGHandler := handler.NewTransactionCategories(db)
	transacttionCTGRoutes := r.Group("/transcat")
	transacttionCTGRoutes.POST("/create", transacttionCTGHandler.Create)
	transacttionCTGRoutes.GET("/read/:id", transacttionCTGHandler.Read)
	transacttionCTGRoutes.PATCH("/update/:id", transacttionCTGHandler.Update)
	transacttionCTGRoutes.DELETE("/delete/:id", transacttionCTGHandler.Delete)
	transacttionCTGRoutes.GET("/list", transacttionCTGHandler.List)

	// grouping route with /transaction
	transactionHandler := handler.NewTransaction(db)
	transactionRoutes := r.Group("/transaction")
	transactionRoutes.POST("/create", middleware.AuthMiddleware(signingKey), transactionHandler.NewTransaction)
	transactionRoutes.GET("/list", middleware.AuthMiddleware(signingKey), transactionHandler.TransactionList)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func NewDatabase() *gorm.DB {
	// dsn := "host=localhost port=5432 user=postgres dbname=digi sslmode=disable TimeZone=Asia/Jakarta"
	dsn := os.Getenv("DATABASE")
	if dsn == "" {
		log.Fatal("DATABASE not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get DB object: %v", err)
	}

	var currentDB string
	err = sqlDB.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		log.Fatalf("failed to query current database: %v", err)
	}

	log.Printf("Current Database: %s\n", currentDB)

	return db
}
