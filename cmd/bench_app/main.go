package main

import (
	"github.com/joho/godotenv"
	"github.com/nikitalystsev/BookSmart/internal/bench_app"
	"log"
)

const configsDir = "configs"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	bench_app.Run(configsDir)
}
