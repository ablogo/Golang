package repository

import (
	"fmt"
	"log/slog"

	"watcher/models"

	_ "modernc.org/sqlite"
)

func (d *Sqlyte) SaveRecord(record models.Record) error {

	query := `INSERT INTO records (app, message, description, is_error, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := d.Db.Exec(query, record.App, record.Message, record.Description, record.IsError, record.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
	}

	return nil
}

func (d *Sqlyte) GetRecords() ([]models.Record, error) {

	query := "SELECT * FROM records"
	rows, err := d.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to query records: %w", err)
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var r models.Record
		err = rows.Scan(&r.Id, &r.App, &r.Message, &r.Description, &r.IsError, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan record: %w", err)
		}
		records = append(records, r)
	}

	return records, nil
}

func (d *Sqlyte) GetLogs() ([]models.Log, error) {

	query := "SELECT * FROM logs"
	rows, err := d.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var r models.Log
		err = rows.Scan(&r.Id, &r.Level, &r.Message, &r.Source, &r.Attributes, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan log: %w", err)
		}
		logs = append(logs, r)
	}

	return logs, nil
}
