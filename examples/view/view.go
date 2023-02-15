package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/kafka/table"
	"github.com/ionos-cloud/streams/store/memory"
	"github.com/ionos-cloud/streams/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/katallaxie/pkg/server"
	"github.com/spf13/cobra"
)

type service struct {
	view view.View[string]

	server.Listener
}

func (s *service) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		app := fiber.New()
		app.Use(logger.New())

		app.Get("/:key", func(c *fiber.Ctx) error {
			v, err := s.view.Get(c.Params("key"))
			if errors.Is(err, view.ErrCatchup) {
				return c.SendStatus(fiber.StatusServiceUnavailable)
			}

			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}

			return c.SendString(v)
		})

		app.Post("/:key", func(c *fiber.Ctx) error {
			err := s.view.Set(c.Params("key"), string(c.Body()))
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			return c.SendStatus(fiber.StatusOK)
		})

		app.Listen(":3000")

		return nil
	}
}

var rootCmd = &cobra.Command{
	Use: "view",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
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
	name := table.Topic("test")
	table := table.WithContext(ctx, table.WithTopic(name), table.WithBrokers("localhost:9092"))

	v := view.New(table, codec.StringEncoder, codec.StringDecoder, store)

	srv := &service{
		view: v,
	}

	s, _ := server.WithContext(ctx)
	s.Listen(v, true)
	s.Listen(srv, false)

	if err := s.Wait(); !errors.Is(err, &server.Error{}) {
		return err
	}

	return nil
}
