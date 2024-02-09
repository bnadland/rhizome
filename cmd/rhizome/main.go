package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bnadland/rhizome/internal/web"
	"github.com/caarlos0/env/v10"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Addr        string `env:"ADDR" envDefault:":3000"`
}

func main() {
	cfg := config{}
	if err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := web.Run(ctx, cfg.Addr, cfg.DatabaseURL); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
