package main

import (
	"log"

	"github.com/davidkamere/htmx-go-scaffolding/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
