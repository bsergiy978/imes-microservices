package main

import (
	"github.com/COVIDEV/viq-chat-services/rooms/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	rooms "github.com/COVIDEV/viq-chat-services/rooms/proto/rooms"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.github.romatroskin.viqchat.rooms.service.rooms"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	rooms.RegisterRoomsHandler(service.Server(), handler.NewHandler(service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
