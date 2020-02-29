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
	teleBot := bot.New(telegram)
	app.Post("/telegram/"+telegram.AccessToken, teleBot.Handler)

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
