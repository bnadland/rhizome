package main

import (
	"log"
	"os"

	"github.com/bnadland/rhizome/internal/web"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(*cli.Context) error {
			return web.Run(":3000")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
