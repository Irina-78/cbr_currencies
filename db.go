package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlCreateTable = `
        CREATE TABLE IF NOT EXISTS cbr_exchange_rate(
            rate_date TEXT NOT NULL,
            num_code INTEGER NOT NULL,
            currency_name TEXT NOT NULL,
            char_code TEXT NOT NULL,
            denomination INTEGER NOT NULL,
            rate_value FLOAT NOT NULL,
            PRIMARY KEY(rate_date, num_code)
        );`

	sqlInsertItem = `
        INSERT OR REPLACE INTO cbr_exchange_rate
            (rate_date, num_code, currency_name, char_code, denomination, rate_value)
            VALUES(?, ?, ?, ?, ?, ?);`
)

type DbStorage struct {
	name string
}

func newDbStorage(name string) *DbStorage {
	return &DbStorage{name: name}
}

// Prepares the database for work.
func (s *DbStorage) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", s.name)
	if err != nil {
		return fmt.Errorf("failed to open the database: %v", err)
	}
	defer db.Close()

	if _, err = s.ExecQuery(ctx, sqlCreateTable); err != nil {
		return fmt.Errorf("failed to create a table: %v", err)
	}

	return nil
}

// Adds a new record in the database.
func (s *DbStorage) Add(ctx context.Context, query *ExchRateQuery, currencies *Currencies,
	filter *CurrencyFilter) error {

	for _, c := range *currencies {
		if filter.IsEnabled() && !filter.IsCurrencyEnabled(c.CharCode) {
			continue
		}

		rows, err := s.ExecQuery(ctx, sqlInsertItem,
			query.Date("2006-01-02"),
			c.NumCode,
			c.Name,
			c.CharCode,
			c.Nominal,
			c.Value,
		)
		if err != nil || rows == 0 {
			return fmt.Errorf("failed to insert a currency: %v", err)
		}
	}

	return nil
}

func (s *DbStorage) ExecQuery(ctx context.Context, query string, params ...any) (int64, error) {
	var (
		db    *sql.DB
		tx    *sql.Tx
		stmt  *sql.Stmt
		res   sql.Result
		count int64
		err   error
	)

	if db, err = sql.Open("sqlite3", s.name); err != nil {
		return 0, fmt.Errorf("failed to open the database: %v", err)
	}
	defer db.Close()

	if tx, err = db.Begin(); err != nil {
		return 0, fmt.Errorf("failed to begin a transaction: %v", err)
	}
	if stmt, err = tx.PrepareContext(ctx, query); err != nil {
		return 0, fmt.Errorf("incorrect query: %v", err)
	}
	defer stmt.Close()

	if res, err = stmt.ExecContext(ctx, params...); err != nil {
		return 0, fmt.Errorf("database query failed: %v", err)
	}
	if count, err = res.RowsAffected(); err != nil {
		return 0, fmt.Errorf("unknown database query execution status: %v", err)
	}
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit a transaction: %v", err)
	}

	return count, nil
}

func (s *DbStorage) SelectCount(ctx context.Context, query string, params ...any) (int, error) {
	var (
		db    *sql.DB
		stmt  *sql.Stmt
		count int
		err   error
	)

	if db, err = sql.Open("sqlite3", s.name); err != nil {
		return 0, fmt.Errorf("failed to open the database: %v", err)
	}
	defer db.Close()

	if stmt, err = db.PrepareContext(ctx, query); err != nil {
		return 0, fmt.Errorf("incorrect query: %v", err)
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, params...).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to get the value from the database: %v", err)
	}

	return count, nil
}
