package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/kafka"
	"github.com/ionos-cloud/streams/store/memory"
	"github.com/ionos-cloud/streams/view"
	"github.com/katallaxie/pkg/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "view",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.SilenceUsage = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	store := memory.New()
	table := kafka.NewTable()

	v := view.New[string](table, streams.StringEncoder{}, streams.StringDecoder{}, store)

	s, _ := server.WithContext(ctx)
	s.Listen(v, false)

	return nil
}
