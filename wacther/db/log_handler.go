package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteHandler struct {
	db   *sql.DB
	opts slog.HandlerOptions
}

// NewSQLiteHandler creates the table if not exist and returns the handler
func NewSQLiteHandler(db *sql.DB, opts slog.HandlerOptions) *SQLiteHandler {
	query := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		level TEXT,
		message TEXT,
		source TEXT,
		attributes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		slog.Error(err.Error())
	}

	return &SQLiteHandler{db: db, opts: opts}
}

func (h *SQLiteHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := h.opts.Level
	if minLevel == nil {
		minLevel = slog.LevelInfo
	}
	return level >= minLevel.Level()
}

// Handle processes the log record and inserts it into SQLite
func (h *SQLiteHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract attributes into a map to convert to JSON
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	attrJSON, err := json.Marshal(attrs)
	if err != nil {
		attrJSON = []byte("{}")
	}

	// Extract source file tracking if requested
	var sourceFile string
	if h.opts.AddSource && r.PC != 0 {
		sourceFile = fmt.Sprintf("%s:%s:%d", r.Source().File, r.Source().Function, r.Source().Line)
	}

	// Execute async or sync database write
	query := `INSERT INTO logs (created_at, level, message, source, attributes) VALUES (?, ?, ?, ?, ?)`
	_, err = h.db.ExecContext(ctx, query,
		r.Time.Format(time.RFC3339),
		r.Level.String(),
		r.Message,
		sourceFile,
		string(attrJSON),
	)
	return err
}

// WithAttrs and WithGroup are required by the interface for creating sub-loggers
func (h *SQLiteHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *SQLiteHandler) WithGroup(name string) slog.Handler       { return h }
