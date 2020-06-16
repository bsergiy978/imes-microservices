package main

import (
	"github.com/COVIDEV/viq-chat-services/jwt/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	jwt "github.com/COVIDEV/viq-chat-services/jwt/proto/jwt"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.github.romatroskin.viqchat.jwt.service.jwt"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	jwt.RegisterJwtHandler(service.Server(), new(handler.Jwt))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
