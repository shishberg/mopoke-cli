package cmd

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Prefix      string             `bson:"prefix"`
	Num         int                `bson:"num"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}
