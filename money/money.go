package money

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	// ErrUnexpectedCSVLine if the csv line does not have the proper format
	ErrUnexpectedCSVLine = errors.New("unexpected format in csv line")
)

// Transaction represents a single transaction from the file
type Transaction struct {
	ID      string
	LabelID string
	Amount  float64
	Account string
	Date    time.Time
	Comment string
}

// NewTransaction creates a new transaction
func NewTransaction(date time.Time, amount float64, account, comment string) (t Transaction, err error) {
	id := generateID(account, date, amount, comment)

	t = Transaction{
		ID:      id,
		LabelID: "1",
		Amount:  amount,
		Account: account,
		Date:    date,
		Comment: comment,
	}

	return t, nil
}

// NewTransactionFromCSV creates a new transaction from a WellsFargoFormatted CSV line
func NewTransactionFromCSV(account string, rawCSV []string) (t Transaction, err error) {
	if len(rawCSV) != 5 {
		return t, ErrUnexpectedCSVLine
	}

	currTime, err := time.Parse("01/02/2006", rawCSV[0])
	if err != nil {
		return t, err
	}

	amount, err := strconv.ParseFloat(rawCSV[1], 64)
	if err != nil {
		return t, err
	}

	comment := rawCSV[4]

	return NewTransaction(currTime, amount, account, comment)
}

// InsertInitialTransaction inserts an initial transaction signifying the initial amount the bank account had on it
func InsertInitialTransaction(transactions []Transaction, account string, amount float64) []Transaction {
	if amount == 0.0 {
		return transactions
	}

	comment := "Initial Balance"
	date := findOldestTransaction(transactions)
	t := Transaction{
		ID:      generateInitialID(account),
		LabelID: "1",
		Amount:  amount,
		Account: account,
		Date:    date,
		Comment: comment,
	}

	return append(transactions, t)
}

func findOldestTransaction(transactions []Transaction) (oldest time.Time) {
	oldest = time.Now()
	for _, t := range transactions {
		if t.Date.Before(oldest) {
			oldest = t.Date
		}
	}
	return oldest
}

// generateID generates a unique id for any transaction you want to insert
func generateID(account string, date time.Time, amount float64, comment string) (id string) {
	idStr := fmt.Sprintf("%s-%s-%f-%s", account, date.Format("01/02/2006"), amount, comment)
	return fmt.Sprintf("%x", md5.Sum([]byte(idStr)))
}

// generateInitialID generates a unique ID that is only to be used for the initial transaction
func generateInitialID(account string) (id string) {
	idStr := fmt.Sprintf("%s-INITIAL-BALANCE", account)
	return fmt.Sprintf("%x", md5.Sum([]byte(idStr)))
}
