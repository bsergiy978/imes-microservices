package main

import (
	otp "github.com/COVIDEV/viq-chat-services/otp/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	proto "github.com/COVIDEV/viq-chat-services/otp/proto/otp"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.github.romatroskin.viqchat.otp.service.otp"),
		micro.Version("latest"),
	)

	// Initialize service
	service.Init()

	// Register Handler
	proto.RegisterOtpHandler(service.Server(), otp.NewOtpHandler(service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
