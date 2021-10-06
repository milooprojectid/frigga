package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"

	bot "frigga/modules/bot"
	driver "frigga/modules/driver"
	http "frigga/modules/httphandler"
	d "frigga/modules/providers/discord"
	m "frigga/modules/providers/messenger"
	service "frigga/modules/service"
)

func newApp() *iris.Application {
	app := iris.Default()

	// TEGERAM HANDLER
	telegram := bot.GetProvider("telegram")
	teleBot := bot.New(telegram)
	app.Post("/telegram/"+os.Getenv("TELEGRAM_SECRET"), teleBot.Handler)
	// ---

	// NAVER LINE HANDLER
	line := bot.GetProvider("line")
	lineBot := bot.New(line)
	app.Post("/line/"+os.Getenv("LINE_SECRET"), lineBot.Handler)
	// ---

	// FACEBOOK MESSEGER HANDLER
	messenger := bot.GetProvider("messenger")
	messengerBot := bot.New(messenger)
	messengerPath := "/messenger/" + os.Getenv("MESSENGER_SECRET")
	app.Post(messengerPath, messengerBot.Handler)
	app.Get(messengerPath, m.VerifySignature)
	// ---

	// DISCORD HANDLER
	discord := bot.GetProvider("discord")
	discordBot := bot.New(discord)
	app.Post("/discord", d.VerifySignature, discordBot.Handler)
	// ---

	// HTTP HANDLERS
	app.Post("subscription/notify", http.AuthGuardMiddleware, http.SendNotificationToSubscriptionHandler)
	app.Post("broadcast", http.AuthGuardMiddleware, http.SendBroadcastMessageHandler)
	app.Post("direct/{sessionId}", http.AuthGuardMiddleware, http.SendDirectMessageHandler)
	app.Post("worker", http.AuthGuardMiddleware, http.BotWorkerHandler)
	app.Get("/", http.HealthCheck)

	// ---
	return app
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading environment file")
	}

	// initialize singletons
	driver.InitializeFirestore()
	driver.InitializeCloudTask()

	service.InitializeServices()
	service.InitializeGrpcServices()

	bot.RegisterCommands()
	// --

	app := newApp()
	port := os.Getenv("APP_PORT")
	app.Run(iris.Addr(":" + port))
}
