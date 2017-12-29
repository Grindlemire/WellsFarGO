package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	bolt "github.com/coreos/bbolt"
	"github.com/grindlemire/seezlog"
	"github.com/jessevdk/go-flags"
)

// Opts are options you can pass into the cli
type Opts struct {
	File   string `short:"f" long:"file" description:"The csv file you want to parse" required:"true"`
	DBPath string `short:"d" long:"db" description:"The DB file you will be persisting in" default:"mydb.bolt"`
}

// Transaction represents a single transaction from the file
type Transaction struct {
	Amount  float64
	Date    time.Time
	Comment string
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

	// Parse incoming file
	transactions, err := parseFile(opts.File)
	if err != nil {
		log.Error("Error parsing file: ", err)
		exit(1)
	}

	// Insert into DB file
	err = persistTransactions(transactions, opts.DBPath)
	if err != nil {
		log.Error("Error persisting transactions: ", err)
		exit(1)
	}

	// Get current budget file

	// Update current budget file
}

func persistTransactions(ts []Transaction, dbPath string) (err error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) (err error) {
		accountBucket := tx.Bucket([]byte("Checking"))
		if accountBucket == nil {
			accountBucket, err = tx.CreateBucket([]byte("Checking"))
			if err != nil {
				return err
			}
		}

		for _, t := range ts {
			timeBucket := accountBucket.Bucket([]byte(t.Date.Format("01/02/2006")))
			if timeBucket == nil {
				timeBucket, err = accountBucket.CreateBucket([]byte(t.Date.Format("01/02/2006")))
				if err != nil && err != bolt.ErrBucketExists {
					return err
				}
			}

			timeBucket.Put([]byte(strconv.FormatFloat(t.Amount, 'f', 2, 64)), []byte(t.Comment))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func parseFile(file string) (transactions []Transaction, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	rawCSV, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, t := range rawCSV {
		if len(t) != 5 {
			return transactions, ErrUnexpectedFileFormat
		}

		currTime, err := time.Parse("01/02/2006", t[0])
		if err != nil {
			return transactions, err
		}

		amount, err := strconv.ParseFloat(t[1], 64)
		if err != nil {
			return transactions, err
		}

		comment := t[4]

		currTransaction := Transaction{
			Amount:  amount,
			Comment: comment,
			Date:    currTime,
		}
		transactions = append(transactions, currTransaction)
	}

	return transactions, nil

}

func exit(status int) {
	log.Flush()
	os.Exit(status)
}
