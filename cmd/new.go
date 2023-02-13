package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shishberg/mopoke-cli/db"

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
	var rels []db.Rel
	for _, rel := range newArgs.relIn {
		r, err := parseRel(rel)
		if err != nil {
			log.Fatal(err)
		}
		rels = append(rels, r)
	}

	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		tickets := client.Database("mopoke").Collection("tickets")
		relationships := client.Database("mopoke").Collection("rel")

		session, err := client.StartSession()
		if err != nil {
			return errors.Trace(err)
		}
		defer session.EndSession(ctx)

		result, err := session.WithTransaction(ctx,
			func(sessCtx mongo.SessionContext) (interface{}, error) {
				ticket := db.Ticket{
					Name:        newArgs.name,
					Title:       newArgs.title,
					Description: newArgs.description,
				}
				insertResult, err := tickets.InsertOne(sessCtx, ticket)
				if err != nil {
					return nil, errors.Trace(err)
				}
				newID := insertResult.InsertedID.(primitive.ObjectID)

				for _, rel := range rels {
					var other db.Ticket
					if err := tickets.FindOne(sessCtx, bson.D{{"name", rel.OtherName}}).Decode(&other); err != nil {
						return nil, errors.Trace(err)
					}
					if rel.ToOther {
						rel.From = newID
						rel.To = other.ID
					} else {
						rel.From = other.ID
						rel.To = newID
					}
					if _, err := relationships.InsertOne(sessCtx, rel); err != nil {
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

func parseRel(r string) (db.Rel, error) {
	delim := strings.IndexAny(r, "<>")
	if delim == -1 {
		return db.Rel{}, errors.New("rel must be of the form type<name or type>name")
	}
	rel := db.Rel{
		Type:      r[:delim],
		OtherName: r[delim+1:],
		ToOther:   delim == '>',
	}
	return rel, nil
}
