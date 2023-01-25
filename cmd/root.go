package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
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
