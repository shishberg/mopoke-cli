package cmd

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"github.com/shishberg/mopoke-cli/db"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	schemaCmd = &cobra.Command{
		Use: "schema",
		Run: runSchema,
	}
)

func init() {
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) {
	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		schemata := client.Database("mopoke").Collection("schemata")

		cursor, err := schemata.Find(ctx, bson.D{{}})
		if err != nil {
			return errors.Trace(err)
		}
		var sch []db.Schema
		if err := cursor.All(ctx, &sch); err != nil {
			return errors.Trace(err)
		}
		for _, s := range sch {
			fmt.Println(s)
		}
		return nil
	})
}
