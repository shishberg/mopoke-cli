package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var (
	pingCmd = &cobra.Command{
		Use: "ping",
		Run: runPing,
	}
)

func runPing(cmd *cobra.Command, args []string) {
	withMongo(cmd, func(ctx context.Context, client *mongo.Client) error {
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return err
		}
		return nil
	})
}
