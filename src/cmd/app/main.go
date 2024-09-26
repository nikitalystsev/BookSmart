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

// @title BookSmart API
// @version 1.0
// @description API Server for BookSmart Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	app.Run(configsDir)

}
