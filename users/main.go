package main

import (
	users "github.com/COVIDEV/viq-chat-services/users/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/COVIDEV/viq-chat-services/users/proto/users"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.github.romatroskin.viqchat.users.service.users"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterUsersHandler(service.Server(), users.NewHandler(service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
