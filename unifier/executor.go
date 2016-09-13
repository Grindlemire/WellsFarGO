package unifier

import (
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/iterator"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/query"
)

// QueryDateRange query for a specific date in the graph
func (u Unifier) QueryDateRange(start, end time.Time) (results []query.Node, err error) {

	path := cayley.StartPath(u.db).
		Filter(iterator.CompareLT, quad.Int(end.Unix())).
		Filter(iterator.CompareGTE, quad.Int(start.Unix())).In()

	nodes, err := path.Iterate(nil).AllValues(nil)
	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryDay query for a specific date in the graph
func (u Unifier) QueryDay(day time.Time) (results []query.Node, err error) {

	path := cayley.StartPath(u.db).
		Has("date", quad.Int(day.Unix()))

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryAmount query for a specific amount
func (u Unifier) QueryAmount(amount float64) (results []query.Node, err error) {

	path := cayley.StartPath(u.db).
		Has("amount", quad.Float(amount))

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

// QueryAmountRange query for a specific range of ammounts
func (u Unifier) QueryAmountRange(lower, upper float64) (results []query.Node, err error) {

	path := cayley.StartPath(u.db).
		Filter(iterator.CompareLT, quad.Float(upper)).
		Filter(iterator.CompareGTE, quad.Float(lower)).In()

	nodes, err := path.Iterate(nil).AllValues(nil)

	results, err = u.parseResults(nodes...)
	return results, err
}

func (u Unifier) parseResults(ids ...quad.Value) (results []query.Node, err error) {
	// QueryNodes query a node by id
	for _, id := range ids {
		m := map[string]interface{}{}

		for _, t := range []string{"date", "amount", "location"} {
			var nodeVals []quad.Value
			nodeVals, err = cayley.StartPath(u.db, id).Out(t).Iterate(nil).AllValues(nil)
			if err != nil {
				log.Error("Error getting value of node: ", err)
				return []query.Node{}, nil
			}
			if len(nodeVals) > 0 {
				m[t] = quad.NativeOf(nodeVals[0])
			}

		}

		uTime := int64(m["date"].(int))
		n := query.Node{
			ID:       quad.StringOf(id),
			Date:     time.Unix(uTime, 0),
			Amount:   m["amount"].(float64),
			Location: m["location"].(string),
		}

		results = append(results, n)
	}
	return results, nil
}
