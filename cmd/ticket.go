package cmd

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`

	Rel      []Rel    `bson:"rel"`      // aggregated
	Children []Ticket `bson:"children"` // aggregated
}

type Rel struct {
	From primitive.ObjectID `bson:"from"`
	To   primitive.ObjectID `bson:"to"`
	Type string             `bson:"type"`

	toOther   bool   `bson:"-"`
	otherName string `bson:"-"`
}
