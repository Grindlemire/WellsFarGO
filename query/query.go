package query

import (
	"reflect"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
)

// RQuery a query of a specific relationship
type RQuery struct {
	Name string
	Val  interface{}
	Type string
}

// Node a node returned from the graph database
type Node struct {
	ID       string
	Date     time.Time
	Location string
	Amount   float64
}

// Query a query wrapper to send to cayley
type Query struct {
	DB        *cayley.Handle
	Nodes     []string
	Relations []RQuery
	results   []Node
}
type parseFunc func(path *path.Path, relation RQuery) (newPath *path.Path)

var mapTypeToParse = map[string]map[string]parseFunc{
	"In": map[string]parseFunc{
		"string":  parseStringInVal,
		"float64": parseFloatInVal,
	},
	"Out": map[string]parseFunc{
		"string":  parseStringOutVal,
		"float64": parseFloatOutVal,
	},
	"Both": map[string]parseFunc{
		"string":  parseStringBothVal,
		"float64": parseFloatBothVal,
	},
}

// Execute executes a query on the graph
func (q Query) Execute() (results []Node, err error) {

	path := cayley.StartPath(q.DB)

	for _, relation := range q.Relations {
		queryType, tErr := getType(relation.Val)
		if tErr != nil {
			log.Error("Error executing query: ", err)
			continue
		}
		path = mapTypeToParse[relation.Type][queryType](path, relation)
	}

	nodes, err := path.Iterate(nil).AllValues(nil)
	if err != nil {
		return []Node{}, log.Error(err)
	}

	results = q.QueryNodes(nodes...)

	return results, err
}

// getType gets the type of the val in order to parse correctly
func getType(val interface{}) (string, error) {
	if val == nil {
		return "", log.Error("Value cannot be null in query")
	}
	return reflect.TypeOf(val).String(), nil
}

// QueryNodes query a node by id
func (q Query) QueryNodes(ids ...quad.Value) (results []Node) {
	for _, id := range ids {
		nodeVals, err := cayley.StartPath(q.DB, id).Out().Iterate(nil).AllValues(nil)
		if err != nil {
			log.Errorf("Error retreiving record with id %v: %v", id, err)
			continue
		}

		n := parseNode(id, nodeVals)
		results = append(results, n)
	}
	return results
}

// parseNode Bad implementaion of parsing the struct out of the cayley response
func parseNode(id quad.Value, nodeVals []quad.Value) (n Node) {
	tStr, err := time.Parse("\"01/02/2006\"", quad.StringOf(nodeVals[0]))
	if err != nil {
		log.Error("Error converting Node query date to time: ", err)
	}
	n = Node{
		ID:       quad.StringOf(id),
		Date:     tStr,
		Amount:   quad.NativeOf(nodeVals[1]).(float64),
		Location: quad.StringOf(nodeVals[2]),
	}
	return n
}
