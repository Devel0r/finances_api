package main

import (
    "finances/internal/app"
    "finances/pkg/config"
    "finances/pkg/logger"
    "log"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load(".env.local")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    cfg := config.LoadConfig()
    logger := logger.NewLogger()

    appInstance, err := app.NewApp(cfg, logger)
    if err != nil {
        logger.Error("Failed to create app instance", map[string]interface{}{"error": err})
        os.Exit(1)
    }

    appInstance.Run("8080")
}