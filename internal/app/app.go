package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"finances/internal/controller"
	"finances/internal/repository"
	"finances/internal/usecase"
	"finances/pkg/config"
	"finances/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type App struct {
	server *gin.Engine
	db     *pgx.Conn
	logger *logger.Logger
}

func NewApp(cfg *config.Config, logger *logger.Logger) (*App, error) {
	db, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	userService := usecase.NewUserService(userRepo, txRepo, db)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	r.POST("/users/:userID/deposit", userController.Deposit)
	r.POST("/users/:fromUserID/transfer/:toUserID", userController.Transfer)
	r.GET("/users/:userID/transactions", userController.GetLastTransactions)

	return &App{
		server: r,
		db:     db,
		logger: logger,
	}, nil
}

func (a *App) Run(port string) {
	a.logger.Info("Starting server", map[string]interface{}{"port": port})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.server.Run(":" + port); err != nil {
			a.logger.Error("Failed to start server", map[string]interface{}{"error": err})
			done <- syscall.SIGINT
		}
	}()

	<-done
	a.logger.Info("Shutting down server", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.db.Close(ctx); err != nil {
		a.logger.Error("Failed to close database connection", map[string]interface{}{"error": err})
	} else {
		a.logger.Info("Database connection closed", nil)
	}
}
