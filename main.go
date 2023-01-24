package main

import (
	"context"
	"flag"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var args struct {
	mongoURL string
}

func main() {
	flag.StringVar(&args.mongoURL, "mongodb", "mongodb://mango:27017", "mongodb address")
	flag.Parse()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(args.mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
