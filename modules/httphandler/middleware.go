package httphandler

import (
	"os"

	"github.com/kataras/iris"
)

// AuthGuardMiddleware ...
func AuthGuardMiddleware(ctx iris.Context) {
	token := ctx.GetHeader("Authorization")

	if token != os.Getenv("APP_BASE_TOKEN") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": "TOKEN_INVALID",
		})
		return
	}

	ctx.Next()
	return
}
