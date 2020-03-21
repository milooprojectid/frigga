package storm

import (
	"log"

	grpc "google.golang.org/grpc"
)

// Client ...
var Client *StormServiceClient

// Init ...
func Init(connectionString string) {
	cc, err := grpc.Dial(connectionString, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}
	// defer cc.Close()
	client := NewStormServiceClient(cc)
	Client = &client
}
