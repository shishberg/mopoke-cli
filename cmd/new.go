package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	newArgs struct {
		name          string
		title         string
		description   string
		relIn, relOut []string
	}

	newCmd = &cobra.Command{
		Use: "new",
		Run: runNew,
	}
)

func init() {
	newCmd.Flags().StringVar(&newArgs.name, "name", "", "ticket name")
	newCmd.Flags().StringVar(&newArgs.title, "title", "", "ticket title")
	newCmd.Flags().StringVar(&newArgs.description, "description", "", "ticket description")
	newCmd.Flags().StringArrayVar(&newArgs.relIn, "rel", nil, "related tickets")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) {
	var rels []Rel
	for _, rel := range newArgs.relIn {
		r, err := parseRel(rel)
		if err != nil {
			log.Fatal(err)
		}
		rels = append(rels, r)
	}

	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		tickets := client.Database("mopoke").Collection("tickets")

		session, err := client.StartSession()
		if err != nil {
			return errors.Trace(err)
		}
		defer session.EndSession(ctx)

		result, err := session.WithTransaction(ctx,
			func(sessCtx mongo.SessionContext) (interface{}, error) {
				ticket := Ticket{
					Name:        newArgs.name,
					Title:       newArgs.title,
					Description: newArgs.description,
					Rel:         []Rel{},
				}
				insertResult, err := tickets.InsertOne(sessCtx, ticket)
				if err != nil {
					return nil, errors.Trace(err)
				}
				id := insertResult.InsertedID.(primitive.ObjectID)

				for _, rel := range rels {
					var other Ticket
					if err := tickets.FindOne(sessCtx, bson.D{{"name", rel.otherName}}).Decode(&other); err != nil {
						return nil, errors.Trace(err)
					}
					out := Rel{
						Type:  rel.Type,
						Other: other.ID,
						Dir:   rel.Dir,
					}
					in := Rel{
						Type:  rel.Type,
						Other: id,
						Dir:   -rel.Dir,
					}
					if _, err := tickets.UpdateByID(sessCtx, id, bson.D{{
						"$push", bson.D{{"rel", out}},
					}}); err != nil {
						return nil, errors.Trace(err)
					}
					if _, err := tickets.UpdateByID(sessCtx, other.ID, bson.D{{
						"$push", bson.D{{"rel", in}},
					}}); err != nil {
						return nil, errors.Trace(err)
					}
				}

				return insertResult, nil
			})
		if err != nil {
			return errors.Trace(err)
		}
		fmt.Println(result)
		return nil
	})
}

func parseRel(r string) (Rel, error) {
	delim := strings.IndexAny(r, "<>:")
	if delim == -1 {
		return Rel{}, errors.New("rel must be of the form type<name, type>name or type:name")
	}
	rel := Rel{
		Type:      r[:delim],
		otherName: r[delim+1:],
	}
	switch r[delim] {
	case '<':
		rel.Dir = -1
	case '>':
		rel.Dir = 1
	case ':':
		rel.Dir = 0
	}
	return rel, nil
}
