package main

import (
	"Booksmart/internal/config"
	"github.com/joho/godotenv"
	"github.com/nikitalystsev/BookSmart-tech-ui/requesters"
	"log"
)

const configDir = "configs"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg, err := config.Init(configDir)
	if err != nil {
		return
	}

	requester := requesters.NewRequester(
		cfg.Auth.JWT.AccessTokenTTL,
		cfg.Auth.JWT.RefreshTokenTTL,
		cfg.Port,
	)

	requester.Run()
}
