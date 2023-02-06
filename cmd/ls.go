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
		cursor, err := tickets.Find(ctx, bson.D{{}})
		if err != nil {
			return err
		}
		var results []Ticket
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}

		for _, t := range results {
			fmt.Printf("%s %s\n", t.Name, t.Title)
		}
		return nil
	})
}
