package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"

	telegram "frigga/modules/telegram"
)

func newApp() *iris.Application {
	app := iris.Default()

	telegramBot := telegram.NewBot(os.Getenv("TELEGRAM_TOKEN"))
	app.Post("/telegram", telegramBot.Handler)

	return app
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading environment file")
	}

	app := newApp()
	port := os.Getenv("APP_PORT")
	app.Run(iris.Addr(":" + port))
}
