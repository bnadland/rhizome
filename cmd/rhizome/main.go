package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bnadland/rhizome/internal/web"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := web.Run(ctx, ":3000"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
