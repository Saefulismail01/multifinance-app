package delivery

import (
	"log"
	"net/http"
	"time"

	"multifinance/config"
	"multifinance/delivery/controller"
	"multifinance/repository"
	"multifinance/usecase/transaction"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Run initializes all dependencies and starts the HTTP server.
func Run() {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	// Connect to database
	sqlxDB, cfg, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer sqlxDB.Close()

	log.Printf("Successfully connected to database %s on %s:%s", cfg.DBName, cfg.Host, cfg.Port)

	// Initialize repositories
	customerRepo := repository.NewCustomerRepository(sqlxDB)
	limitRepo := repository.NewLimitRepository(sqlxDB)
	transactionRepo := repository.NewTransactionRepository(sqlxDB)

	// Initialize usecase
	transactionUsecase := transaction.NewTransactionUsecase(
		transactionRepo,
		customerRepo,
		limitRepo,
		sqlxDB,
	)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		transactionHandler := controller.NewTransactionHandler(transactionUsecase)
		transactionHandler.RegisterRoutes(v1)
	}

	// Start the server
	serverAddr := ":" + cfg.APIConfig.ApiPort
	log.Printf("Server is running on http://localhost%s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
