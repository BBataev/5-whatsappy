package postgres

import (
	"fmt"
	"log/slog"

	"github.com/BBataev/whatsappy/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func BuildDSN(cfg *config.Config) {
	connStr := fmt.Sprintf(
		"user=%s host=%s port=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgreUser,
		cfg.PostgreHost,
		cfg.PostgrePort,
		cfg.PostgrePass,
		cfg.PostgreName,
	)

	var err error

	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
	}
}

func CloseCon() {
	if db != nil {
		_ = db.Close()
	}
}
