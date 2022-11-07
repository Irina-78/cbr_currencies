package main

import (
	"context"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilename = "test.db"

const sqlSelect = `
	SELECT COUNT(*)
        FROM cbr_exchange_rate
        WHERE rate_date = ?
            AND num_code = ?
            AND currency_name = ?
            AND char_code = ?
            AND denomination = ?
            AND rate_value = ?;`

func TestDbStorage(t *testing.T) {
	var err error

	ctx := context.Background()

	storage := newDbStorage(dbFilename)
	if err = storage.Init(ctx); err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	query := newExchRateQuery()

	cs := Currencies{
		Currency{
			NumCode:  840,
			CharCode: "USD",
			Nominal:  1,
			Name:     "US Dollar",
			Value:    59.9756,
		},
		Currency{
			NumCode:  978,
			CharCode: "EUR",
			Nominal:  1,
			Name:     "Euro",
			Value:    62.5903,
		},
	}

	if err = storage.Add(ctx, query, &cs, newCurrencyFilter()); err != nil {
		t.Fatalf("failed to insert data: %v", err)
	}

	var count int
	for _, c := range cs {
		count, err = storage.SelectCount(
			ctx,
			sqlSelect,
			query.Date("2006-01-02"),
			c.NumCode,
			c.Name,
			c.CharCode,
			c.Nominal,
			c.Value,
		)
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		if count != 1 {
			t.Fatalf("expected 1 got %d", count)
		}
	}
}

func TestDeleteDbFile(t *testing.T) {
	if err := os.Remove(dbFilename); err != nil {
		t.Fatalf("failed to delete database file: %v", err)
	}
}
