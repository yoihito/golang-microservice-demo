package main

import (
	"auth/config"
	"auth/internal/app"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(config)
}
