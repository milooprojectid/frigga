package morbius

import (
	"log"

	grpc "google.golang.org/grpc"
)

// Client ...
var Client *MorbiusServiceClient

// Init ...
func Init(connectionString string) {
	cc, err := grpc.Dial(connectionString, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}
	// defer cc.Close()
	client := NewMorbiusServiceClient(cc)
	Client = &client
}
