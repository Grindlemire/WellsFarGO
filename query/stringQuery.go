package query

import (
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/cihub/seelog"
)

func parseStringInVal(path *path.Path, relation RQuery) (newPath *path.Path) {
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
			newPath = path.Has(relation.Name, quad.String(relation.Val.(string)))
		}
	}
	return newPath
}

func parseStringOutVal(path *path.Path, relation RQuery) (newPath *path.Path) {
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
			newPath = path.Has(relation.Name, quad.String(relation.Val.(string)))
		}
	}
	return newPath
}

func parseStringBothVal(path *path.Path, relation RQuery) (newPath *path.Path) {
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
			newPath = path.Has(relation.Name, quad.String(relation.Val.(string)))
		}
	}
	return newPath
}
