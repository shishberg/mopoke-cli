package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoURL string
	timeout  time.Duration

	rootCmd = &cobra.Command{
		Use: "mopoke",
		Run: run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&mongoURL, "db", "mongodb://localhost:27017", "mongodb address")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "timeout")
}

func run(cmd *cobra.Command, args []string) {
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
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
