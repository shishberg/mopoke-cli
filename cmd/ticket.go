package cmd

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Rel         []Rel              `bson:"rel"`
}

type Rel struct {
	Type  string             `bson:"type"`
	Other primitive.ObjectID `bson:"other"`
	Dir   int                `bson:"dir"`

	otherName string `bson:"-"`
}
