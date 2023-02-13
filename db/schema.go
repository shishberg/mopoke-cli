package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	ID    primitive.ObjectID    `bson:"_id,omitempty"`
	Types map[string]TypeSchema `bson:"inline"`
}

type TypeSchema struct {
	Fields map[string]string `bson:"inline"`
}

func (s Schema) String() string {
	var str string
	for t, ts := range s.Types {
		str += fmt.Sprintf("%s:\n", t)
		for f, ft := range ts.Fields {
			str += fmt.Sprintf("  %s: %s\n", f, ft)
		}
	}
	return str
}
