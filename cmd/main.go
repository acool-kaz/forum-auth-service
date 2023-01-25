package main

import (
	"log"

	"github.com/acool-kaz/forum-auth-service/internal/app"
	"github.com/acool-kaz/forum-auth-service/internal/config"
)

const configPath = "./config.json"

func main() {
	cfg, err := config.InitConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.InitApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.RunApp()
}
