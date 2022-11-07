package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	argCurrency []string
	argDate     []string
	argSql      string
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cbr_currencies",
		Short: "Gets the Bank of Russia exchange rate",
		Long:  "cbr_currencies is a tool to get the Bank of Russia exchange rate for today or specified date.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isArgsEmpty(args) {
				logger.Info(fmt.Sprintf("args was entered: %v", args))
				return fmt.Errorf("unknown command was entered: %q",
					strings.Join(args, ", "))
			}

			if len(argCurrency) > 0 {
				logger.Info(fmt.Sprintf("currencies was entered: %v", argCurrency))
				for i, c := range argCurrency {
					c = strings.ToUpper(strings.TrimSpace(c))
					if currencyFilter.CodeExists(c) {
						argCurrency[i] = c
					} else {
						return fmt.Errorf("currency value %q is incorrect", c)
					}
				}
			}

			if len(argDate) > 0 {
				logger.Info(fmt.Sprintf("dates was entered: %v", argDate))
				for i, d := range argDate {
					d = strings.TrimSpace(d)
					if isDateCorrect(d) {
						argDate[i] = d
					} else {
						return fmt.Errorf("date value %q is incorrect", d)
					}
				}
			}

			if len(argSql) > 0 {
				logger.Info(fmt.Sprintf("database file name was entered: %s", argSql))
				argSql = strings.TrimSpace(argSql)
				if !isFileNameCorrect(argSql) {
					return fmt.Errorf("invalid database file name: %q", argSql)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&argCurrency, "currency", "c", []string{},
		"the currency you're interested, for example 'USD' (according to ISO 4217)")
	cmd.Flags().StringSliceVarP(&argDate, "date", "d", []string{},
		"exchange rate date (as day.month.year)")
	cmd.Flags().StringVarP(&argSql, "sql", "s", "",
		"name of the SQLite database file in which the exchange rate data should be saved, for example 'currencies.db'")
	cmd.Flags().SortFlags = false

	return cmd
}

// Checks the entered arguments are empty
func isArgsEmpty(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, a := range args {
		if strings.TrimSpace(a) != "" {
			return false
		}
	}
	return true
}

// Checks the entered date.
func isDateCorrect(s string) bool {
	// default date is 2 Jan 2006
	if _, err := time.Parse("2.1.06", s); err != nil {
		if _, err = time.Parse("2.1.2006", s); err != nil {
			return false
		}
	}
	return true
}

// Checks the entered file name.
func isFileNameCorrect(s string) bool {
	re := regexp.MustCompile(`^[A-Za-zА-Яа-я]{1,}[A-Za-zА-Яа-я0-9.]{0,}$`)
	return re.MatchString(s)
}
