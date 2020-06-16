package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"

	messages "github.com/COVIDEV/viq-chat-services/messages/proto/messages"
	rooms "github.com/COVIDEV/viq-chat-services/rooms/proto/rooms"
	users "github.com/COVIDEV/viq-chat-services/users/proto/users"
	"github.com/centrifugal/centrifuge"
)

func handleLog(e centrifuge.LogEntry) {
	log.Infof("%s: %v", e.Message, e.Fields)
}

func waitExitSignal(n *centrifuge.Node) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		n.Shutdown(context.Background())
		done <- true
	}()
	<-done
}

func main() {
	cfg := centrifuge.DefaultConfig

	// Set HMAC secret to handle requests with JWT auth too. This is
	// not required if you don't use token authentication and
	// private subscriptions verified by token.
	cfg.TokenHMACSecretKey = "9884c0c0-5a98-401c-ad0f-591f7001278f"
	cfg.LogLevel = centrifuge.LogLevelDebug
	cfg.LogHandler = handleLog

	cfg.JoinLeave = true
	cfg.HistoryLifetime = 300
	cfg.HistorySize = 1000
	cfg.HistoryRecover = true

	cfg.UserSubscribeToPersonal = true
	cfg.UserPersonalChannelNamespace = "service"

	cfg.Namespaces = []centrifuge.ChannelNamespace{
		{
			Name: "chat",
			ChannelOptions: centrifuge.ChannelOptions{
				Publish:            true,
				SubscribeToPublish: true,
				Presence:           true,
				JoinLeave:          true,
			},
		},
		{
			Name: "service",
			ChannelOptions: centrifuge.ChannelOptions{
				Publish:            true,
				SubscribeToPublish: true,
				Presence:           true,
				JoinLeave:          true,
			},
		},
	}

	service := web.NewService(
		web.Name("com.github.romatroskin.viqchat.centrifuge.service"),
		web.Version("latest"),
		web.Address(":8000"),
	)

	// Initialise service
	service.Init()

	_ = users.NewUsersService("com.github.romatroskin.viqchat.users.service.users", client.DefaultClient)
	roomsService := rooms.NewRoomsService("com.github.romatroskin.viqchat.rooms.service.rooms", client.DefaultClient)
	messagesService := messages.NewChatMessagesService("com.github.romatroskin.viqchat.messages.service.messages", client.DefaultClient)

	node, _ := centrifuge.New(cfg)

	node.On().ClientConnecting(func(ctx context.Context, t centrifuge.TransportInfo, e centrifuge.ConnectEvent) centrifuge.ConnectReply {
		return centrifuge.ConnectReply{
			Data: []byte(`{}`),
		}
	})

	/// Service messages structure would be like:
	/// {
	///   "type": 0,
	///   "data": {
	///     ...
	///   }
	/// }
	/// where type is one of:
	/// ROOMS = 0 | MESSAGES = 1 | ADD = 2 | REMOVE = 3,
	/// for example:
	///
	/// Initial service message
	/// {
	///   "type": 0,
	///   "data": {
	///     "rooms": []
	///   }
	/// }
	///
	/// Add service message
	/// {
	///   "type": 2,
	///   "data": {
	///     "from": "<userId>",
	///		"room": "<roomId>"
	///   }
	/// }
	///
	/// Remove service message
	/// {
	///   "type": 2,
	///   "data": {
	///		"room": "<roomId>"
	///   }
	/// }
	///

	node.On().ClientConnected(func(ctx context.Context, client *centrifuge.Client) {

		chats, err := roomsService.List(ctx, &rooms.ListRequest{UserId: client.UserID()})
		if err != nil {
			log.Error(err.Error())
		}

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		chatsJSON, _ := json.MarshalToString(chats)
		_, err = node.Publish(node.PersonalChannel(client.UserID()), []byte(`{"type":0,"data":`+chatsJSON+`}`))
		if err != nil {
			log.Error(err.Error())
		}

		client.On().Subscribe(func(e centrifuge.SubscribeEvent) centrifuge.SubscribeReply {
			channel := strings.Split(e.Channel, ":")
			if channel[0] == "chat" {
				if channel[1][:1] == "#" {
					ids := strings.Split(channel[1][1:], ",")
					log.Infof("user %s wants to chat with %s", ids[0], ids[1])
					newRoom, err := roomsService.Create(ctx, &rooms.CreateRequest{Room: &rooms.Room{Topic: e.Channel, Owner: ids[0]}, Participants: ids})
					if err != nil {
						log.Error(err.Error())
					}

					messages, err := messagesService.List(ctx, &messages.ListRequest{Channel: e.Channel})
					if err != nil {
						log.Error(err.Error())
					}

					var json = jsoniter.ConfigCompatibleWithStandardLibrary
					messagesJSON, _ := json.MarshalToString(messages)
					_, err = node.Publish(node.PersonalChannel(client.UserID()), []byte(`{"type":1,"data":`+messagesJSON+`}`))
					if err != nil {
						log.Error(err.Error())
					}

					roomJSON, _ := json.MarshalToString(newRoom)
					_, err = node.Publish(node.PersonalChannel(ids[1]), []byte(`{"type":2,"data":`+roomJSON+`}`))
					if err != nil {
						log.Error(err.Error())
					}
				} else {
					log.Infof("user %s subscribes on %s", client.UserID(), e.Channel)
				}
			}
			return centrifuge.SubscribeReply{}
		})

		client.On().Unsubscribe(func(e centrifuge.UnsubscribeEvent) centrifuge.UnsubscribeReply {
			log.Infof("user %s unsubscribed from %s", client.UserID(), e.Channel)
			return centrifuge.UnsubscribeReply{}
		})

		/// Message formats as follows:
		/// { "from": "userId", "type": 0, "content": "blah blah blah", "parent": "null" }
		/// Look at the proto  for further details.
		client.On().Publish(func(e centrifuge.PublishEvent) centrifuge.PublishReply {
			channel := strings.Split(e.Channel, ":")
			if channel[0] == "chat" {
				if channel[1][:1] == "#" {
					ids := strings.Split(channel[1][1:], ",")
					log.Infof("user %s sent private message to %s", ids[0], ids[1])

					message := &messages.ChatMessage{}
					var json = jsoniter.ConfigCompatibleWithStandardLibrary
					err := json.Unmarshal(e.Data, message)
					if err != nil {
						log.Error(err.Error())
					}

					_, err = messagesService.Write(ctx, message)
					if err != nil {
						log.Error(err.Error())
					}
				} else {
					log.Infof("user %s publishes into channel %s: %s", client.UserID(), e.Channel, string(e.Data))
				}
			}
			return centrifuge.PublishReply{}
		})

		client.On().Message(func(e centrifuge.MessageEvent) centrifuge.MessageReply {
			log.Infof("message from user: %s, data: %s", client.UserID(), string(e.Data))
			return centrifuge.MessageReply{}
		})

		client.On().Disconnect(func(e centrifuge.DisconnectEvent) centrifuge.DisconnectReply {
			log.Infof("user %s disconnected, disconnect: %s", client.UserID(), e.Disconnect)
			return centrifuge.DisconnectReply{}
		})

		transport := client.Transport()
		log.Infof("user %s connected via %s with protocol: %s", client.UserID(), transport.Name(), transport.Protocol())
	})

	node.On().ClientRefresh(func(ctx context.Context, client *centrifuge.Client, e centrifuge.RefreshEvent) centrifuge.RefreshReply {
		log.Infof("user %s connection is going to expire, refreshing", client.UserID())
		return centrifuge.RefreshReply{
			ExpireAt: time.Now().Unix() + 10,
		}
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}

	service.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{}))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	waitExitSignal(node)
}
