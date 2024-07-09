package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	bot "github.com/shlau/supremacy-day-change/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot.BotToken = os.Getenv("BOT_TOKEN")
	bot.Done = make(chan bool)
	bot.Run()
}
