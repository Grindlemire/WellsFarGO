package unifier

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"time"

	"github.com/cayleygraph/cayley"
	cgraph "github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt" // Needed to use bolt as the backend of cayley
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/format"
	"github.com/grindlemire/WellsFarGO/query"
	// "github.com/cayleygraph/cayley/quad"
)

const dbExistsErr = "quadstore: cannot init; database already exists"

// Unifier will read data and unify it with the db
type Unifier struct {
	db      *cayley.Handle
	newData format.Transactions
}

// NewUnifier creates a unifier object that will read data and unify it with the db
func NewUnifier(dbFile, csvFile, formatType string) (u *Unifier, err error) {
	err = cgraph.InitQuadStore("bolt", dbFile, nil)
	if err != nil && err.Error() != dbExistsErr {
		log.Error("Error initializing cayley quad store: ", err)
		return nil, err
	}

	graph, err := cayley.NewGraph("bolt", dbFile, nil)
	if err != nil {
		log.Error("Error creating cayley graph database: ", err)
	}

	f := format.FormatMap[formatType]
	newData, err := format.ReadFiles(csvFile, f)
	if err != nil {
		return nil, err
	}

	u = &Unifier{
		db:      graph,
		newData: newData,
	}

	return u, err

}

// AddNewData adds the new parsed data to the cayley database
func (u Unifier) AddNewData() (err error) {
	for time, transactions := range u.newData {
		for _, transaction := range transactions {
			err = u.AddTransaction(time, transaction)
			if err != nil {
				log.Error("Error inserting transaction into cayley: ", err)
				continue
			}
		}
	}
	return err
}

// AddTransaction adds a transaction to the cayley database
func (u Unifier) AddTransaction(t time.Time, trans format.Transaction) (err error) {
	timeStr := t.Format("01/02/2006")
	tranByte, err := json.Marshal(trans)
	if err != nil {
		return log.Error("Error marshalling transaction: ", err)
	}
	timeByte := []byte(timeStr)
	key := append(tranByte, timeByte...)

	hasher := md5.New()
	hasher.Write(key)
	idStr := hex.EncodeToString(hasher.Sum(nil))

	err = u.AddRelationship(idStr, "date", timeStr)
	if err != nil {
		return err
	}

	err = u.AddRelationship(idStr, "amount", trans.Amount)
	if err != nil {
		rErr := u.RemoveRelationship(idStr, "date", timeStr)
		if rErr != nil {
			log.Critical("Error Rolling back after an error: ", rErr)
		}
		return err
	}

	err = u.AddRelationship(idStr, "location", trans.Location)
	if err != nil {
		rErr := u.RemoveRelationship(idStr, "date", timeStr)
		if rErr != nil {
			log.Critical("Error Rolling back after an error: ", rErr)
		}
		rErr = u.RemoveRelationship(idStr, "amount", trans.Amount)
		if rErr != nil {
			log.Critical("Error Rolling back after an error: ", rErr)
		}
		return err
	}

	return err
}

// AddRelationship Adds a relationship to the graph db
func (u Unifier) AddRelationship(n, r string, val interface{}) (err error) {
	if !checkType(val) {
		return log.Error("Must store a primitive as a value")
	}

	q := &query.Query{
		DB: u.db,
	}
	results := q.QueryNodes(quad.String(n))
	if len(results) > 0 {
		return nil
	}

	err = u.db.AddQuad(quad.Make(n, r, val, nil))
	if err != nil {
		log.Error("Error inserting relationship: ", err)
	}
	return err
}

// RemoveRelationship Removes a relationship from the graph db
func (u Unifier) RemoveRelationship(n, r string, val interface{}) (err error) {
	if !checkType(val) {
		return log.Error("Must be a primitive as value")
	}

	err = u.db.RemoveQuad(quad.Make(n, r, val, nil))
	if err != nil {
		log.Error("Error removing relationship: ", err)
	}
	return err
}

// QueryDate query for a specific date in the graph
func (u Unifier) QueryDate(t time.Time) (results []query.Node, err error) {

	q := &query.Query{
		DB:    u.db,
		Nodes: []string{},
		Relations: []query.RQuery{
			query.RQuery{
				Name: "date",
				Val:  t.Format("01/02/2006"),
				Type: "Out",
			},
		},
	}

	results, err = q.Execute()
	if err != nil {
		log.Error("Error Executing Query: ", err)
		return []query.Node{}, nil
	}

	return results, err
}

func checkType(val interface{}) bool {
	theType := reflect.TypeOf(val).String()
	switch theType {
	case "string":
		return true
	case "int":
		return true
	case "float64":
		return true
	}
	return false
}
