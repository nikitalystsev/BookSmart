package main

import (
	"Booksmart/internal/app"
	"github.com/joho/godotenv"
	"log"
)

const configsDir = "configs"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	app.Run(configsDir)

}
