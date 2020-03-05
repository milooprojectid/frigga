package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"

	bot "frigga/modules/bot"
	driver "frigga/modules/driver"
	service "frigga/modules/service"
)

func newApp() *iris.Application {
	app := iris.Default()

	telegram := bot.GetProvider("telegram")
	line := bot.GetProvider("line")
	messenger := bot.GetProvider("line")

	teleBot := bot.New(telegram)
	lineBot := bot.New(line)
	messengerBot := bot.New(messenger)

	app.Post("/telegram/"+os.Getenv("TELEGRAM_SECRET"), teleBot.Handler)
	app.Post("/line/"+os.Getenv("LINE_SECRET"), lineBot.Handler)
	app.Post("/messenger/"+os.Getenv("MESSENGER_SECRET"), messengerBot.Handler)

	return app
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading environment file")
	}

	// initialize singletons
	driver.InitializeFirestore()
	service.InitializeServices()
	bot.RegisterCommands()

	app := newApp()
	port := os.Getenv("APP_PORT")
	app.Run(iris.Addr(":" + port))
}
