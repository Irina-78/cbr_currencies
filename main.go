package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

const maxInstances = 3

var (
	logger         *zap.Logger // all methods are safe for concurrent use
	currencyFilter *CurrencyFilter
)

func init() {
	var err error
	logger, err = newProductionLogger()
	if err != nil {
		fmt.Printf("can't initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	currencyFilter = newCurrencyFilter()
}

func main() {
	var err error

	cmd := newRootCmd()
	if err = cmd.Execute(); err != nil {
		// cobra prints an error
		logger.Info(fmt.Sprintf("cmd error: %v", err))
		return
	}
	if cmd.Flags().Changed("help") {
		// cobra prints help
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cmd.Flags().Changed("currency") {
		if len(argCurrency) == 0 {
			logger.Warn("entered an empty currency")

			fmt.Println("pass the currency you're interested, ",
				"for example \"-c USD\".")
			return
		}

		for _, c := range argCurrency {
			if err = currencyFilter.CurrencyEnable(c); err != nil {
				logger.Warn(fmt.Sprintf("currency %q wasn't enabled: %v", c, err))

				fmt.Printf("%v.\nPass the currency you're interested, for example \"-c USD\".\n", err)
				return
			}
		}

		currencyFilter.Enable()
	}

	var queries []*ExchRateQuery
	if cmd.Flags().Changed("date") {
		len := len(argDate)
		if len == 0 {
			logger.Warn("entered an empty date")

			fmt.Println("pass the exchange rate date you're interested, ",
				"for example \"-d 01.01.2022\".")
			return
		}

		queries = make([]*ExchRateQuery, len, len)
		for i, d := range argDate {
			q := newExchRateQuery()
			if err = q.SetDate(d); err != nil {
				logger.Warn(fmt.Sprintf("setting date %q failed: %v", d, err))

				fmt.Printf("%v.\nPass the exchange rate date you're interested as day.month.year.\n", err)
				return
			}
			queries[i] = q
		}
	} else {
		queries = []*ExchRateQuery{newExchRateQuery()}
	}

	var storage *DbStorage
	if cmd.Flags().Changed("sql") {
		if len(argSql) == 0 {
			logger.Warn("entered an empty name of the database file")

			fmt.Println("pass the name of the database file in which the exchange rate data ",
				"should be saved, for example \"-s currencies.db\".")
			return
		}

		storage = newDbStorage("./" + argSql)
		if err = storage.Init(ctx); err != nil {
			logger.Error(fmt.Sprintf("failed to create the database: %v", err))

			fmt.Printf("failed to create the database: %v\n", err)
			return
		}
	}

	var wg sync.WaitGroup

	printer := newResultPrinter()

	for _, query := range queries {
		for runtime.NumGoroutine() > maxInstances {
			time.Sleep(500 * time.Millisecond)
		}

		wg.Add(1)
		go worker(&wg, ctx, query, currencyFilter, printer, storage)
	}

	wg.Wait()

	fmt.Println("Done.")
}

func worker(wg *sync.WaitGroup, ctx context.Context, query *ExchRateQuery,
	filter *CurrencyFilter, printer *ResultPrinter, storage *DbStorage) {

	defer wg.Done()

	client := newCbrClient()
	answer, err := client.Get(ctx, query.String())
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] failed: %v", query, err))

		fmt.Printf("request %q wasn't completed: %v\n", query, err)
		return
	}
	logger.Debug(fmt.Sprintf("[%s] received an answer: %s", query, answer))

	decoder := newCbrDecoder(answer)
	if !decoder.isValid {
		logger.Error(fmt.Sprintf("[%s] failed, received incorrect answer: %s", query, answer))

		fmt.Printf("response to request %q was not decoded, received incorrect answer:\n%s\n", query, answer)
		return
	}

	result := CbrResult{}
	if err = decoder.Decode(&result); err != nil {
		logger.Error(fmt.Sprintf("[%s] decoding failed: %v", query, err))

		fmt.Printf("response to request %q was not decoded: %v\n", query, err)
		return
	}
	logger.Info(fmt.Sprintf("[%s] response successfully decoded", query))

	// print the answer
	printer.print(query, filter, &result.Currencies)

	// save the answer to db
	if storage != nil {
		err = storage.Add(ctx, query, &result.Currencies, filter)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to save data to the database: %v", err))

			fmt.Printf("failed to save data to the database: %v\n", err)
			return
		}
		logger.Info(fmt.Sprintf("[%s] data successfully saved in %q", query, storage.name))
	}
}
