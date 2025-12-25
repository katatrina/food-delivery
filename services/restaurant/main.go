package main

import (
	"os"

	"github.com/katatrina/food-delivery/services/restaurant/cmd/app"
)

func main() {
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
