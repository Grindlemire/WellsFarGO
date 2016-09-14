package unifier

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/cayleygraph/cayley"
	cgraph "github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt" // Needed to use bolt as the backend of cayley
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/format"
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
	unixTime := t.Unix()
	tranByte, err := json.Marshal(trans)
	if err != nil {
		return log.Error("Error marshalling transaction: ", err)
	}
	timeByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeByte, uint64(unixTime))
	key := append(tranByte, timeByte...)

	hasher := md5.New()
	hasher.Write(key)
	idStr := hex.EncodeToString(hasher.Sum(nil))

	if u.NodeExists(quad.String(idStr)) {
		return nil
	}

	err = u.AddRelationship(idStr, "date", unixTime)
	if err != nil {
		return err
	}

	err = u.AddRelationship(idStr, "amount", trans.Amount)
	if err != nil {
		rErr := u.RemoveRelationship(idStr, "date", unixTime)
		if rErr != nil {
			log.Critical("Error Rolling back after an error: ", rErr)
		}
		return err
	}

	err = u.AddRelationship(idStr, "location", trans.Location)
	if err != nil {
		rErr := u.RemoveRelationship(idStr, "date", unixTime)
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

func checkType(val interface{}) bool {
	switch val.(type) {
	case string:
		return true
	case int64, int32:
		return true
	case float64, float32:
		return true
	}
	return false
}
