package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var (
	pingCmd = &cobra.Command{
		Use: "ping",
		Run: ping,
	}
)

func ping(cmd *cobra.Command, args []string) {
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
	fmt.Println("pong")
}
