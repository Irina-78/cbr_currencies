package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ExchRateQuery struct {
	link string
	time time.Time
}

// Creates an 'ExchRateQuery' instance.
func newExchRateQuery() *ExchRateQuery {
	return &ExchRateQuery{
		link: "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=",
		time: time.Now(),
	}
}

// Returns the set date in the given format according to the Time.Format specification.
func (q *ExchRateQuery) Date(format string) string {
	return q.time.Format(format)
}

// Takes a string in format "day.month.year" and sets the date in the query.
func (q *ExchRateQuery) SetDate(date string) error {
	var dt time.Time
	var err error
	date = strings.TrimSpace(date)

	// default date is 2 Jan 2006
	if dt, err = time.Parse("2.1.06", date); err != nil {
		if dt, err = time.Parse("2.1.2006", date); err != nil {
			return fmt.Errorf("incorrect date format: %v", err)
		}
	}

	q.time = dt

	return nil
}

// Builds the query string.
func (q *ExchRateQuery) String() string {
	var s strings.Builder
	s.WriteString(q.link)
	s.WriteString(q.time.Format("02/01/2006"))
	return s.String()
}

type Currency struct {
	NumCode  int
	CharCode string
	Nominal  int
	Name     string
	Value    float64
}

type Currencies []Currency

type CbrResult struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies Currencies `xml:"Valute"`
}

func (c Currency) String() string {
	return fmt.Sprintf("%8d %s\t%10.4f RUB", c.Nominal, c.CharCode, c.Value)
}

// 'CurrencyFilter' filters interested currencies, if any have been set.
type CurrencyFilter struct {
	enabled bool
	list    map[string]bool
}

// Creates a new disabled `CurrencyFilter` instance.
func newCurrencyFilter() *CurrencyFilter {
	return &CurrencyFilter{
		enabled: false,
		list: map[string]bool{
			"AUD": false,
			"AZN": false,
			"AMD": false,
			"BYN": false,
			"BGN": false,
			"BRL": false,
			"HUF": false,
			"KRW": false,
			"HKD": false,
			"DKK": false,
			"USD": false,
			"EUR": false,
			"INR": false,
			"KZT": false,
			"CAD": false,
			"KGS": false,
			"CNY": false,
			"MDL": false,
			"TMT": false,
			"NOK": false,
			"PLN": false,
			"RON": false,
			"XDR": false,
			"SGD": false,
			"TJS": false,
			"TRY": false,
			"UZS": false,
			"UAH": false,
			"GBP": false,
			"CZK": false,
			"SEK": false,
			"CHF": false,
			"ZAR": false,
			"JPY": false,
		},
	}
}

// Checks the correctness of the currency code.
func (f *CurrencyFilter) CodeExists(code string) bool {
	_, exist := f.list[code]
	return exist
}

func (f *CurrencyFilter) IsEnabled() bool {
	return f.enabled
}

func (f *CurrencyFilter) Enable() {
	f.enabled = true
}

func (f *CurrencyFilter) Disable() {
	f.enabled = false
}

func (f *CurrencyFilter) IsCurrencyEnabled(code string) bool {
	if _, ok := f.list[code]; ok {
		return f.list[code]
	}
	return false
}

func (f *CurrencyFilter) IsCurrencyDisabled(code string) bool {
	return !f.IsCurrencyEnabled(code)
}

func (f *CurrencyFilter) CurrencyEnable(code string) error {
	if _, ok := f.list[code]; ok {
		f.list[code] = true
		return nil
	}
	return fmt.Errorf("currency code is incorrect: %s", code)
}

// Http client.
type CbrClient struct {
	client *http.Client
}

// Creates a 'CbrClient' instance.
func newCbrClient() *CbrClient {
	tr := &http.Transport{
		IdleConnTimeout:   5 * time.Second,
		DisableKeepAlives: true,
	}

	return &CbrClient{
		client: &http.Client{Transport: tr},
	}
}

// Makes a request to the server to get exchange rate data.
func (c *CbrClient) Get(ctx context.Context, query string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", query, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request '%s': %v", query, err)
	}
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0`)
	req.Header.Add("Connection", "close")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("request '%s' failed: %v", query, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}

type ResultPrinter struct {
	sync.Mutex
}

func newResultPrinter() *ResultPrinter {
	return &ResultPrinter{}
}

func (w *ResultPrinter) print(query *ExchRateQuery, filter *CurrencyFilter, currencies *Currencies) {
	w.Lock()
	defer w.Unlock()

	fmt.Printf("\nData on %s\n", query.Date("02.01.2006"))
	for _, c := range *currencies {
		if filter.IsEnabled() && filter.IsCurrencyDisabled(c.CharCode) {
			continue
		}
		fmt.Println(c)
	}
}
