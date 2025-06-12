package main

import (
	"log/slog"
	"os"

	"github.com/BBataev/whatsappy/internal/config"
	"github.com/BBataev/whatsappy/internal/service/endpoints"
	"github.com/BBataev/whatsappy/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	cfg     *config.Config
)

func main() {
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Logger Initialized")

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file")
	}

	cfg = config.Load()

	slog.Info("Config loaded", slog.Any("config", cfg))

	postgres.BuildDSN(cfg)
	slog.Info("Connected to postgres")

	defer func() {
		postgres.CloseCon()
		slog.Info("Connection to postgres closed")
	}()

	s := gin.Default()

	s.POST("/register", endpoints.HandleRegister)
	s.POST("/login", endpoints.HandleLogin)

	s.Run(cfg.ListenAddr)

}
