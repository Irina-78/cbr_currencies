package main

import (
	"testing"
)

func TestCurrencyCodeExists(t *testing.T) {
	f := newCurrencyFilter()
	if f.CodeExists("") {
		t.Fatalf("empty string found in the list of currency codes")
	}
	if f.CodeExists("RUB") {
		t.Fatalf("'RUB' found in the list of currency codes")
	}

	for c, _ := range f.list {
		if !f.CodeExists(c) {
			t.Fatalf("currency code '%s' doesn't exist", c)
		}
	}
}

func TestNewCurrencyFilterIsDisabled(t *testing.T) {
	if newCurrencyFilter().IsEnabled() {
		t.Fatalf("new currency filter is enabled")
	}
}

func TestCurrencyFilterIsEnabled(t *testing.T) {
	f := newCurrencyFilter()

	f.Enable()
	if !f.IsEnabled() {
		t.Fatalf("expected true got %v", f.IsEnabled())
	}

	f.Disable()
	if f.IsEnabled() {
		t.Fatalf("expected false got %v", f.IsEnabled())
	}
}

func TestCurrencyFilterEnableCodes(t *testing.T) {
	var err error
	currencies := []string{"USD", "EUR", "AUD"}
	f := newCurrencyFilter()

	for _, c := range currencies {
		if err = f.CurrencyEnable(c); err != nil {
			t.Fatalf("currency code '%s' not found", c)
		}
		if !f.IsCurrencyEnabled(c) {
			t.Fatalf("currency code '%s' isn't enabled", c)
		}
	}

	code := "RUB"
	if err = f.CurrencyEnable(code); err == nil {
		t.Fatalf("currency code '%s' found", code)
	}
	if f.IsCurrencyEnabled(code) {
		t.Fatalf("currency code '%s' is enabled", code)
	}
}

func TestExchRateQuery(t *testing.T) {
	var err error
	q := newExchRateQuery()
	if err = q.SetDate("122020"); err == nil {
		t.Fatalf("incorrect date format")
	}

	q = newExchRateQuery()
	if err = q.SetDate(" 01.12.2020 "); err != nil {
		t.Fatalf("failed to set the date: %v", err)
	}
	expected := "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=01/12/2020"
	if q.String() != expected {
		t.Fatalf("expected %s got %s", expected, q.String())
	}
}
