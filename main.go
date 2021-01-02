package main

import (
	"log"
	"os"

	"dont-genshin-bot/app"
	"github.com/joho/godotenv"
)

const envFile = "./config.env"

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	if _, err := os.Stat(envFile); err == nil {
		_ = godotenv.Load(envFile)
	}
}

func main() {
	app.Start()
}
