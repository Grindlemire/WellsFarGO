package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/grindlemire/WellsFarGO/money"
	"github.com/grindlemire/WellsFarGO/postgres"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/seezlog"
	"github.com/jessevdk/go-flags"
)

// Opts are options you can pass into the cli
type Opts struct {
	Account       string   `short:"a" long:"account" description:"The account that this file is for" default:"checking"`
	Files         []string `short:"f" long:"file" description:"The csv file you want to parse" required:"true"`
	InitialAmount float64  `long:"initial-amount" description:"The initial amount of money in the bank account" default:"0"`
}

var (
	// ErrUnexpectedFileFormat is thrown when the file is a csv but of the wrong format compared to what WellsFargo puts out
	ErrUnexpectedFileFormat = errors.New("unexpected file format")
)

var opts Opts
var parser = flags.NewParser(&opts, flags.Default)

func main() {
	logger, err := seezlog.SetupConsoleLogger(seezlog.Info)
	if err != nil {
		fmt.Printf("Error setting up logger: %v\n", err)
		exit(1)
	}
	err = log.ReplaceLogger(logger)
	if err != nil {
		fmt.Printf("Error replacing logger: %v\n", err)
		exit(1)
	}
	defer log.Flush()

	_, err = parser.Parse()
	if err != nil {
		return
	}

	log.Infof("Config: %#v", opts)

	var transactions []money.Transaction
	for _, file := range opts.Files {
		// Parse incoming file
		currTransactions, err := parseFile(opts.Account, file)
		if err != nil {
			log.Error("Error parsing file: ", err)
			exit(1)
		}

		transactions = append(transactions, currTransactions...)
	}

	transactions = money.InsertInitialTransaction(transactions, opts.Account, opts.InitialAmount)

	c, err := postgres.NewConnection()
	if err != nil {
		log.Error("Error establishing connection to postgres: ", err)
		exit(1)
	}

	rowsInserted, err := c.InsertTransactions(transactions)
	if err != nil {
		log.Error("Error persisting transactions: ", err)
		exit(1)
	}

	log.Infof("Inserted %d new transactions", rowsInserted)
}

func parseFile(account, file string) (transactions []money.Transaction, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	allRawCSV, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, rawCSV := range allRawCSV {
		t, err := money.NewTransactionFromCSV(account, rawCSV)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func exit(status int) {
	log.Flush()
	os.Exit(status)
}
