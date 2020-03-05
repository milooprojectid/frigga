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
	app.Post("/telegram/"+os.Getenv("TELEGRAM_SECRET"), teleBot.Handler)

	line := bot.GetProvider("line")
	lineBot := bot.New(line)
	app.Post("/line/"+os.Getenv("LINE_SECRET"), lineBot.Handler)

	messenger := bot.GetProvider("messenger")
	messengerBot := bot.New(messenger)
	messengerPath := "/messenger/" + os.Getenv("MESSENGER_SECRET")
	app.Post(messengerPath, messengerBot.Handler)
	app.Get(messengerPath, func(ctx iris.Context) {
		secret := os.Getenv("MESSENGER_SECRET")
		mode := ctx.URLParam("hub.mode")
		token := ctx.URLParam("hub.verify_token")
		challenge := ctx.URLParam("hub.challenge")

		if mode == "subscribe" && token == secret {
			ctx.StatusCode(200)
			ctx.Text(challenge)
		} else {
			ctx.StatusCode(403)
			ctx.Text("cant validate signature")
		}
	})

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
