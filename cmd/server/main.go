package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"keeper/internal/core/config"
	"keeper/internal/core/encryption"
	"keeper/internal/core/service"
	"keeper/internal/infra/repo"
	"keeper/internal/infra/s3"
	"keeper/internal/logger"
	"keeper/migrations"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.Server.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	err = migrations.RunMigration(ctx, cfg.Server.DatabaseDSN)
	if err != nil {
		logger.Log.Fatal("Migration error", zap.String("error", err.Error()))
	}

	if err := run(ctx, cfg); err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	logger.Log.Debug("Connecting to DB")
	db, err := sql.Open("pgx", cfg.Server.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize Database: %w", err)
	}
	defer db.Close()

	cancelCtx, cancel := context.WithCancel(ctx)

	logger.Log.Debug("Connecting Repos")
	repoUser := repo.NewUserRepo(db)
	repoSecret := repo.NewSecretRepo(db)

	logger.Log.Debug("Connecting Master Key")
	masterKey, err := encryption.LoadPrivateKey([]byte(cfg.Server.MasterKey))
	if err != nil {
		logger.Log.Fatal("Master Key error", zap.String("error", err.Error()))
	}

	logger.Log.Debug("Creatin EncryptionService")
	encrypter := encryption.NewEncryptionService(cfg.Server.EncryptionKeyPath, masterKey)

	logger.Log.Debug("Creating S3 Client")
	s3Client, err := s3.NewS3Client(ctx, &cfg.S3)
	if err != nil {
		logger.Log.Fatal("Creating S3 client error", zap.String("error", err.Error()))
	}

	iamService := service.NewIAMService(repoUser, &cfg.IAM)
	secretService := service.NewSecretService(repoSecret, s3Client, encrypter, cfg.Server.EncryptionKeyVersion)
	logger.Log.Info("Service initialized")

	fmt.Println(iamService, cancelCtx, secretService) // todo
	// api := rest.NewAPI(cfg, iamService, ordersService, balanceService)

	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#manually
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		fmt.Println("Runned")
		// if err := api.Run(cfg.Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		// 	logger.Log.Info("Runing server error", zap.Error(err))
		// }
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	// sdCtx, cancelAPI := context.WithTimeout(ctx, 5*time.Second)
	// defer cancelAPI()
	// if err := api.Shutdown(sdCtx); err != nil {
	// 	log.Fatal("Server forced to shutdown:", err)
	// }

	logger.Log.Info("Waitng for processing goroutines to finish")
	cancel()
	// procWg.Wait()

	logger.Log.Info("Server exiting")
	return nil
}
