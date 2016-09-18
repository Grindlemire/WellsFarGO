package unifier

import (
	"strings"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/iterator"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
)

// Node a node returned from the graph database
type Node struct {
	ID       string  `json:"id"`
	Date     string  `json:"date"`
	Location string  `json:"location"`
	Amount   float64 `json:"amount"`
}

// QueryDateRange query for a specific date in the graph
func (u Unifier) QueryDateRange(start, end time.Time) (results []Node, err error) {

	path := cayley.StartPath(u.db).
		Filter(iterator.CompareLT, quad.Int(end.Unix())).
		Filter(iterator.CompareGTE, quad.Int(start.Unix())).In()

	nodes, err := path.Iterate(nil).AllValues(nil)
	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryDay query for a specific date in the graph
func (u Unifier) QueryDay(day time.Time) (results []Node, err error) {

	path := cayley.StartPath(u.db).
		Has("date", quad.Int(day.Unix()))

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryAmount query for a specific amount
func (u Unifier) QueryAmount(amount float64) (results []Node, err error) {

	path := cayley.StartPath(u.db).
		Has("amount", quad.Float(amount))

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryAmountRange query for a specific range of ammounts
func (u Unifier) QueryAmountRange(lower, upper float64) (results []Node, err error) {

	path := cayley.StartPath(u.db).
		Filter(iterator.CompareLT, quad.Float(upper)).
		Filter(iterator.CompareGTE, quad.Float(lower)).In()

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryLocation query for a specific location
func (u Unifier) QueryLocation(location string) (results []Node, err error) {

	path := cayley.StartPath(u.db).
		Has("location", quad.String(location))

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

func (u Unifier) parseResults(ids ...quad.Value) (results []Node, err error) {
	// QueryNodes query a node by id
	for _, id := range ids {
		m := map[string]interface{}{}

		for _, t := range []string{"date", "amount", "location"} {
			var nodeVals []quad.Value
			nodeVals, err = cayley.StartPath(u.db, id).Out(t).Iterate(nil).AllValues(nil)
			if err != nil {
				log.Error("Error getting value of node: ", err)
				return []Node{}, nil
			}
			if len(nodeVals) > 0 {
				m[t] = quad.NativeOf(nodeVals[0])
			}

		}

		uTime := int64(m["date"].(int))
		t := time.Unix(uTime, 0)
		tStr := t.Format("01/02/2006")
		idWithQuotes := quad.StringOf(id)
		idArr := strings.Split(idWithQuotes, "\"")
		idStr := idArr[1]

		n := Node{
			ID:       idStr,
			Date:     tStr,
			Amount:   m["amount"].(float64),
			Location: m["location"].(string),
		}

		results = append(results, n)
	}
	return results, nil
}

// NodeExists checks the existence of a node
func (u Unifier) NodeExists(id quad.Value) (exists bool) {
	nodes, _ := cayley.StartPath(u.db, id).Out().Iterate(nil).AllValues(nil)
	if len(nodes) > 0 {
		return true
	}
	return false
}
