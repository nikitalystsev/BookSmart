package main

import (
	"github.com/joho/godotenv"
	"github.com/nikitalystsev/BookSmart/internal/app"
	"log"
)

const configsDir = "configs"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	app.RunEcho(configsDir)
}
