package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	lsCmd = &cobra.Command{
		Use: "ls",
		Run: runLs,
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func runLs(cmd *cobra.Command, args []string) {
	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		tickets := client.Database("mopoke").Collection("tickets")

		var pipeline mongo.Pipeline
		if len(args) > 0 {
			pipeline = append(pipeline, bson.D{{
				"$match", bson.D{{
					"name", bson.M{"$in": args},
				}},
			}})
		}
		pipeline = append(pipeline,
			bson.D{{
				"$graphLookup", bson.M{
					"from":             "rel",
					"startWith":        "$_id",
					"connectFromField": "to",
					"connectToField":   "from",
					"as":               "rel",
					"maxDepth":         0,
				},
			}},
			bson.D{{"$lookup", bson.M{
				"from":         "tickets",
				"localField":   "rel.to",
				"foreignField": "_id",
				"as":           "children",
			}},
			})
		cursor, err := tickets.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		var results []Ticket
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}

		for _, t := range results {
			fmt.Printf("%s %s\n", t.Name, t.Title)
			for _, c := range t.Children {
				fmt.Printf("  %s %s\n", c.Name, c.Title)
			}
		}
		return nil
	})
}
