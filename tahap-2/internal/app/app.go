package app

import (
	"context"
	"log"

	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/handlers"
	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/middleware"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/auth"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/config"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/database"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/worker"
	"mnc-fullstack-technical-test/tahap-2/internal/usecase"
	repo "mnc-fullstack-technical-test/tahap-2/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "mnc-fullstack-technical-test/tahap-2/docs"
)

func Run() {
	cfg := config.LoadConfig()

	// Run SQL migrations (replaces db.AutoMigrate)
	if err := database.RunMigrations(cfg, "tahap-2/migrations"); err != nil {
		log.Printf("Warning: migration error: %v", err)
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		panic(err)
	}

	// Infrastructure
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// Repositories
	userRepo := repo.NewUserRepository(db)
	txRepo := repo.NewTransactionRepository(db)

	// Transfer worker — credit function closes over repos
	creditFn := makeCreditFn(txRepo, userRepo)
	transferWorker := worker.NewTransferWorker(3, 100, creditFn)
	go transferWorker.Start(context.Background())

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)
	txUsecase := usecase.NewTransactionUsecase(txRepo, userRepo, transferWorker)

	// Handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	txHandler := handlers.NewTransactionHandler(txUsecase)
	dashboardHandler := handlers.NewDashboardHandler(transferWorker)

	// Gin Router
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Dashboard (no auth required — monitoring only)
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/queue-stats", dashboardHandler.GetQueueStats)
		dashboard.GET("/jobs", dashboardHandler.GetJobs)
	}

	// Protected Routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.PUT("/profile", userHandler.UpdateProfile)
		protected.POST("/topup", txHandler.TopUp)
		protected.POST("/pay", txHandler.Payment)
		protected.POST("/transfer", txHandler.Transfer)
		protected.GET("/transactions", txHandler.GetTransactions)
	}

	r.Run(":" + cfg.ServerPort)
}

// makeCreditFn returns a worker.CreditFn that updates the receiver balance and
// transaction status using the injected repositories.
func makeCreditFn(txRepo repository.TransactionRepository, userRepo repository.UserRepository) worker.CreditFn {
	return func(ctx context.Context, receiverID uuid.UUID, amount float64, transactionID uuid.UUID, status string) error {
		if err := userRepo.UpdateBalance(ctx, receiverID, amount); err != nil {
			return err
		}
		return txRepo.UpdateStatus(ctx, transactionID, status)
	}
}
