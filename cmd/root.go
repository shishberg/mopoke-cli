package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURL string
	timeout  time.Duration

	rootCmd = &cobra.Command{
		Use: "mopoke",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&mongoURL, "db", "mongodb://localhost:27017", "mongodb address")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "timeout")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func withMongo(cmd *cobra.Command, do func(context.Context, *mongo.Client) error) {
	ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	if err := do(ctx, client); err != nil {
		log.Fatal(err)
	}
}
