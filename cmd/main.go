package main

import (
	"github.com/sayrus383/btc-currency/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
