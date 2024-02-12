package main

import (
	"context"
	"log"
	"os"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	q := db.New(conn)

	if err = q.SavePage(ctx, db.SavePageParams{
		Title: "Home",
		Slug:  "home",
		Content: `# Home
		
hello, [[world]]
		
#public #programming
		`,
	}); err != nil {
		log.Fatal(err)
	}
}
