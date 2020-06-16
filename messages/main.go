package main

import (
	"github.com/COVIDEV/viq-chat-services/messages/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	messages "github.com/COVIDEV/viq-chat-services/messages/proto/messages"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.github.romatroskin.viqchat.messages.service.messages"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	messages.RegisterChatMessagesHandler(service.Server(), handler.NewHandler(service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
