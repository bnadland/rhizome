package main

import (
	"context"
	"log"
	"os"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(*cli.Context) error {
			ctx := context.Background()
			conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
			if err != nil {
				return err
			}
			defer conn.Close(ctx)
			q := db.New(conn)

			return q.SavePage(ctx, db.SavePageParams{
				Title:   "Home",
				Slug:    "home",
				Content: "hello, world",
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
