package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"product-tracking/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBHandler struct {
	DB       *pgxpool.Pool
	Settings *config.Settings
}

func (d *DBHandler) Initialize() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.Settings.DB_USER, d.Settings.DB_PASSWORD, d.Settings.DB_HOST, d.Settings.DB_PORT, d.Settings.DB_NAME)
	ctx := context.Background()
	count := 0

	// Initialize the pgx connection pool
Retry:
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		slog.Error("Unable to connect to database.", "description", err.Error())
		return err
	}

	// Ping the database to ensure connection is live
	if err := pool.Ping(ctx); err != nil {
		count++
		if count <= 2 {
			time.Sleep(time.Second * 20)
			goto Retry
		}

		slog.Error("Database ping failed. Unable to connect to database.")
		return err
	}
	slog.Info("Successfully connected to PostgreSQL with pgxpool.")
	d.DB = pool

	return nil
}

func (s *DBHandler) Close() {
	s.DB.Close()
	slog.Info("Database successfully closed.")
}
