package main

import (
	"github.com/gin-gonic/gin"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/cockroach"
	"github.com/micro/go-micro/v2/web"

	"github.com/COVIDEV/viq-chat-services/auth/api/handler"
	jwt "github.com/COVIDEV/viq-chat-services/jwt/proto/jwt"
	otp "github.com/COVIDEV/viq-chat-services/otp/proto/otp"
	users "github.com/COVIDEV/viq-chat-services/users/proto/users"
)

func main() {
	gin.ForceConsoleColor()

	// New Service
	service := web.NewService(
		web.Name("com.github.romatroskin.viqchat.auth.api"),
		web.Version("latest"),
		web.Address(":6666"),
	)

	// Initialise service
	service.Init()

	otpService := otp.NewOtpService("com.github.romatroskin.viqchat.otp.service.otp", client.DefaultClient)
	jwtService := jwt.NewJwtService("com.github.romatroskin.viqchat.jwt.service.jwt", client.DefaultClient)
	usersService := users.NewUsersService("com.github.romatroskin.viqchat.users.service.users", client.DefaultClient)

	// Create RESTful handler (using Gin)
	api := &handler.Handler{Otp: otpService, Jwt: jwtService, Users: usersService, Store: cockroach.NewStore(
		store.Database("viqchat"),
		store.Table("access"),
	)}

	router := gin.Default()
	router.POST("/api/auth", api.Auth)
	router.POST("/api/verify", api.Verify)
	router.GET("/api/profile", api.Profile)
	router.GET("/api/contacts", api.Contacts)

	// Register Handler
	service.Handle("/", router)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
