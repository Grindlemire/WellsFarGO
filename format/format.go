package format

import (
	"bufio"
	"encoding/csv"
	"os"
	"time"

	log "github.com/cihub/seelog"
)

// Transaction is a record of a purchase
type Transaction struct {
	Location string  `json:"location"`
	Amount   float64 `json:"amount"`
}

// Transactions is a map of a Date to a series of transactions
type Transactions map[time.Time][]Transaction

// FormatMap maps a format to a formatter
var FormatMap = map[string]Formatter{
	"Wells Fargo": WellsFargoFormatter{},
}

// Formatter is an interface for a new format of csv (from other banks) to be added
type Formatter interface {
	ParseTransaction(line []string) (t Transaction, err error)
	ParseTime(line []string) (t time.Time, err error)
}

// ReadFiles reads Well Fargo csv files
func ReadFiles(csvFile string, f Formatter) (newData Transactions, err error) {
	newData = make(Transactions)
	fd, err := os.Open(csvFile)
	if err != nil {
		return nil, log.Error("Error opening csv: ", err)
	}
	r := csv.NewReader(bufio.NewReader(fd))
	fileData, err := r.ReadAll()
	if err != nil {
		return nil, log.Error("Error reading csv: ", err)
	}

	for _, line := range fileData {
		t, tErr := f.ParseTime(line)
		if tErr != nil {
			log.Errorf("Error parsing csv line time: %v. Skipping line.", tErr)
			continue
		}
		trans, tErr := f.ParseTransaction(line)
		if tErr != nil {
			log.Errorf("Error parsing csv line Transaction: %v. Skipping line.", tErr)
			continue
		}
		newData[t] = append(newData[t], trans)
	}

	return newData, err
}
