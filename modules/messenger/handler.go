package messenger

import "github.com/kataras/iris"

// EventAdapter ...
func EventAdapter(ctx iris.Context) ([]TextMessage, error) {
	var message TextMessage
	if err := ctx.ReadJSON(&message); err != nil {
		return []TextMessage{
			message,
		}, err
	}
	return []TextMessage{
		message,
	}, nil
}
