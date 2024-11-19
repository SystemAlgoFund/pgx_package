package pgx_package

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
)

// Database представляет соединение с базой данных.
type Database struct {
	Conn   *pgx.Conn
	Logger *log.Logger
}

// NewDatabase создаёт новое подключение к базе данных.
func NewDatabase(connString string, logger *log.Logger) (*Database, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Database{
		Conn:   conn,
		Logger: logger,
	}, nil
}

// Close закрывает соединение с базой данных.
func (db *Database) Close() error {
	err := db.Conn.Close(context.Background())
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

// Exec выполняет запрос без возврата результата.
func (db *Database) Exec(query string, args ...interface{}) (pgconn.CommandTag, error) {
	db.Logger.Printf("Executing query: %s", query)
	tag, err := db.Conn.Exec(context.Background(), query, args...)
	if err != nil {
		db.Logger.Printf("Query execution failed: %v", err)
		return tag, fmt.Errorf("failed to execute query: %w", err)
	}
	return tag, nil
}

// Query выполняет запрос и сканирует результаты.
func (db *Database) Query(dest interface{}, query string, args ...interface{}) error {
	db.Logger.Printf("Querying: %s", query)
	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		db.Logger.Printf("Query failed: %v", err)
		return fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(dest, rows); err != nil {
		db.Logger.Printf("Scanning failed: %v", err)
		return fmt.Errorf("failed to scan rows: %w", err)
	}
	return nil
}

// Transaction предоставляет возможность выполнения операций в рамках транзакции.
func (db *Database) Transaction(ctx context.Context, txFunc func(pgx.Tx) error) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // перекатить панику
		} else if err != nil {
			tx.Rollback(ctx) // ошибка, откатываем транзакцию
		} else {
			err = tx.Commit(ctx) // коммитим транзакцию
		}
	}()

	err = txFunc(tx)
	return err
}
