package handler

import (
	"context"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/cockroach"

	messages "github.com/COVIDEV/viq-chat-services/messages/proto/messages"
)

type ChatMessages struct {
	store store.Store
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *ChatMessages {
	// Return the initialised store
	return &ChatMessages{
		store: cockroach.NewStore(
			store.Database("viqchat"),
			store.Table("messages"),
		),
	}
}

// Write is an event to store the message to the database
func (e *ChatMessages) Write(ctx context.Context, req *messages.ChatMessage, rsp *messages.Response) error {
	log.Infof("Received message %s request", req)
	// Encode the room
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(req)
	if err != nil {
		return errors.InternalServerError("dev.viqchat.messages.service", "Could not marshal message: %v", err)
	}

	// Write to the store
	if err := e.store.Write(&store.Record{Key: strings.Join([]string{req.Id, req.Channel}, ":"), Value: bytes, Metadata: cockroach.Metadata{
		"from":    req.From,
		"channel": req.Channel,
	}}); err != nil {
		return errors.InternalServerError("dev.viqchat.messages.service", "Could not write to store: %v", err)
	}

	return nil
}

func (e *ChatMessages) List(ctx context.Context, req *messages.ListRequest, rsp *messages.Response) error {
	recs, err := e.store.Read(req.Channel, store.ReadSuffix())
	if err != nil {
		return errors.InternalServerError("dev.viqchat.messages.service", "Could not read from store: %v", err)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	msgs := make([]*messages.ChatMessage, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &msgs[i]); err != nil {
			return errors.InternalServerError("dev.viqchat.rooms.service", "Could not unmarshal room: %v", err)
		}
	}

	rsp.Messages = msgs
	return nil
}
