package query

import (
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
)

func parseFloatInVal(path *path.Path, relation RQuery) (newPath *path.Path) {
	newPath = path

	if relation.Name == "" {
		if relation.Val == "" {
			newPath = path.In()
		} else {
			log.Error("Relation is empty but value is not: ", relation.Val)
		}
	} else {
		if relation.Val == "" {
			newPath = path.In(relation.Name)
		} else {
			newPath = path.Has(relation.Name, quad.Float(relation.Val.(float64))).In()
		}
	}
	return newPath
}

func parseFloatOutVal(path *path.Path, relation RQuery) (newPath *path.Path) {
	newPath = path

	if relation.Name == "" {
		if relation.Val == "" {
			newPath = path.Out()
		} else {
			log.Error("Relation is empty but value is not: ", relation.Val)
		}
	} else {
		if relation.Val == "" {
			newPath = path.Out(relation.Name)
		} else {
			newPath = path.Has(relation.Name, quad.Float(relation.Val.(float64))).Out()
		}
	}
	return newPath
}

func parseFloatBothVal(path *path.Path, relation RQuery) (newPath *path.Path) {
	newPath = path

	if relation.Name == "" {
		if relation.Val == "" {
			newPath = path.Both()
		} else {
			log.Error("Relation is empty but value is not: ", relation.Val)
		}
	} else {
		if relation.Val == "" {
			newPath = path.Both(relation.Name)
		} else {
			newPath = path.Has(relation.Name, quad.Float(relation.Val.(float64))).Both()
		}
	}
	return newPath
}
