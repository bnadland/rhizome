package main

import (
	"log"

	"github.com/bnadland/rhizome/internal/web"
)

func main() {
	if err := web.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
