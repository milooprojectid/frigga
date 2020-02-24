package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
)

func hello(ctx iris.Context) {
	ctx.JSON(map[string]string{
		"message": "hello",
	})
}

func newApp() *iris.Application {
	app := iris.Default()

	app.Post("/line", hello)

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
