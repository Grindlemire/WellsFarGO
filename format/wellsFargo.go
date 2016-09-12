package format

import (
	"strconv"
	"time"
)

// WellsFargoFormatter is a formatter for Wells Fargo csv files
type WellsFargoFormatter struct{}

// ParseTime parses a time from a line in a Wells Fargo csv
func (w WellsFargoFormatter) ParseTime(line []string) (t time.Time, err error) {
	return time.Parse("01/02/2006", line[0])
}

// ParseTransaction parses a transaction from a line in a Wells Fargo csv
func (w WellsFargoFormatter) ParseTransaction(line []string) (t Transaction, err error) {
	amount, err := strconv.ParseFloat(line[1], 64)
	amount *= -1 // Flip so expenses are positive and payments are negative
	if err != nil {
		return Transaction{}, err
	}
	t = Transaction{
		Location: line[4],
		Amount:   amount,
	}
	return t, err
}
