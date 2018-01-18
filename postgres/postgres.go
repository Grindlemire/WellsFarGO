package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	// For the SQL driver for postgres
	_ "github.com/lib/pq"

	"github.com/grindlemire/WellsFarGO/money"
)

// Connection manages connections to postgres
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new connection to postgres
func NewConnection() (c *Connection, err error) {
	connStr, err := createConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	c = &Connection{
		DB: db,
	}

	return c, nil
}

// InsertTransactions will insert transactions to the database
func (c Connection) InsertTransactions(ts []money.Transaction) (rowsInserted int64, err error) {
	for _, t := range ts {
		result, err := c.DB.Exec(
			fmt.Sprintf(`
				INSERT INTO budget.transactions(id, account, labelid, datetime, amount, comment) 
				VALUES ('%s', '%s', '%s', '%s', '%f', $$%s$$)
				ON CONFLICT DO NOTHING`,
				t.ID, t.Account, t.LabelID, t.Date.Format("01/02/2006"), t.Amount, t.Comment),
		)
		if err != nil {
			return 0, err
		}

		batchInserted, err := result.RowsAffected()
		if err != nil {
			return rowsInserted, err
		}

		rowsInserted += batchInserted
	}

	return rowsInserted, nil
}

func createConnectionString() (connStr string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	passFile := fmt.Sprintf("%s/.pgpass.conf", usr.HomeDir)
	f, err := os.Open(passFile)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	creds := strings.Split(string(b), ":")
	if len(creds) != 5 {
		return "", fmt.Errorf("Error parsing credentials file")
	}

	connStr = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", creds[0], creds[1], creds[2], creds[3], creds[4])
	return connStr, err
}
