package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mungkiice/-loan-service/internal/config"
	"github.com/mungkiice/-loan-service/internal/delivery/http"
	"github.com/mungkiice/-loan-service/internal/infrastructure/email"
	"github.com/mungkiice/-loan-service/internal/infrastructure/jwt"
	"github.com/mungkiice/-loan-service/internal/infrastructure/redis"
	"github.com/mungkiice/-loan-service/internal/infrastructure/storage"
	"github.com/mungkiice/-loan-service/internal/repository/postgres"
	"github.com/mungkiice/-loan-service/internal/usecase"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	redisClient, err := redis.NewClient(cfg.Redis.Addr)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	fileStorage, err := storage.NewLocalFileStorage(cfg.Storage.BasePath, cfg.Storage.BaseURL)
	if err != nil {
		log.Fatalf("failed to init file storage: %v", err)
	}

	emailService := email.NewMockEmailService()
	if cfg.Email.Provider == "smtp" {
		// TODO: implement SMTP
		emailService = email.NewMockEmailService()
	}

	loanRepo := postgres.NewLoanRepository(db)
	approvalRepo := postgres.NewApprovalRepository(db)
	investmentRepo := postgres.NewInvestmentRepository(db)
	disbursementRepo := postgres.NewDisbursementRepository(db)
	userRepo := postgres.NewUserRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)
	investorRepo := postgres.NewInvestorRepository(db)

	jwtService := jwt.NewJWTService(cfg.App.JWTSecret, cfg.App.JWTExpiration)

	loanUseCase := usecase.NewLoanUseCase(
		loanRepo,
		approvalRepo,
		investmentRepo,
		disbursementRepo,
		userRepo,
		redisClient,
		fileStorage,
		emailService,
	)

	authUseCase := usecase.NewAuthUseCase(userRepo, employeeRepo, investorRepo, jwtService)

	handler := http.NewHandler(loanUseCase)
	authHandler := http.NewAuthHandler(authUseCase)
	router := http.SetupRouter(handler, authHandler, authUseCase)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	go router.Run(addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
