package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	prefix      string
	title       string
	description string

	newCmd = &cobra.Command{
		Use: "new",
		Run: runNew,
	}
)

func init() {
	newCmd.Flags().StringVar(&prefix, "prefix", "MOP", "ticket ID prefix")
	newCmd.Flags().StringVar(&title, "title", "", "ticket title")
	newCmd.Flags().StringVar(&description, "description", "", "ticket description")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) {
	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		tickets := client.Database("mopoke").Collection("tickets")

		session, err := client.StartSession()
		if err != nil {
			return err
		}
		defer session.EndSession(ctx)

		result, err := session.WithTransaction(ctx,
			func(sessCtx mongo.SessionContext) (interface{}, error) {
				var max Ticket
				opts := options.FindOne().SetSort(bson.D{{"num", -1}})
				maxResult := tickets.FindOne(sessCtx, bson.D{{"prefix", prefix}}, opts)
				if err := maxResult.Decode(&max); err != nil {
					// Not found?
					log.Println(err)
				}
				ticket := Ticket{
					Prefix:      prefix,
					Num:         max.Num + 1,
					Title:       title,
					Description: description,
				}
				fmt.Printf("%#v\n", ticket)
				insertResult, err := tickets.InsertOne(sessCtx, ticket)
				if err != nil {
					return nil, err
				}
				return insertResult, nil
			})
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	})
}
