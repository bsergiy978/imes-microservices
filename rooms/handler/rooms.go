package handler

import (
	"context"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/cockroach"
	"github.com/segmentio/ksuid"

	rooms "github.com/COVIDEV/viq-chat-services/rooms/proto/rooms"
)

// Rooms service
type Rooms struct {
	store     store.Store
	publisher micro.Publisher
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Rooms {
	// Return the initialised store
	return &Rooms{
		store: cockroach.NewStore(
			store.Database("viqchat"),
			store.Table("rooms"),
		),
		publisher: micro.NewPublisher(srv.Name(), srv.Client()),
	}
}

// Create is used to create a room  and store information into the database
func (h *Rooms) Create(ctx context.Context, req *rooms.CreateRequest, rsp *rooms.Response) error {
	log.Infof("Received %s request", req)
	// Validate request
	if req.Room.Topic == "" || req.Room.Owner == "" {
		return errors.BadRequest("dev.viqchat.rooms.service", "Required parameters missing")
	}

	if resp, err := h.findRoomByTopic(req.Room.Topic); err == nil {
		// return errors.Conflict("dev.viqchat.users.service", "User already exist")
		rsp.Room = resp
		return nil
	} else {
		log.Error(err.Error())
	}

	// Add the auto-generated fields
	var room rooms.Room = *req.Room
	if len(room.Id) == 0 {
		// generate ID
		room.Id = ksuid.New().String()
	}

	room.Created = time.Now().Unix()
	room.Updated = time.Now().Unix()

	// Encode the room
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(room)
	if err != nil {
		return errors.InternalServerError("dev.viqchat.rooms.service", "Could not marshal room: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: room.Id, Value: bytes, Metadata: cockroach.Metadata{
		"topic": room.Topic,
	}}); err != nil {
		return errors.InternalServerError("dev.viqchat.rooms.service", "Could not write to store: %v", err)
	}

	go h.publisher.Publish(ctx, &rooms.Event{
		Type: rooms.EventType_RoomCreated,
	})

	for _, v := range req.Participants {
		if err = h.store.Write(&store.Record{Key: strings.Join([]string{v, room.Id}, ":")}, store.WriteTo("viqchat", "user_room")); err != nil {
			return errors.InternalServerError("dev.viqchat.rooms.service", "Could not write to store: %v", err)
		}
	}

	rsp.Room = &room
	return nil
}

// retreives a room given a topic
func (h *Rooms) findRoomByTopic(topic string) (*rooms.Room, error) {
	// Validate request
	if len(topic) == 0 {
		return nil, errors.BadRequest("dev.viqchat.rooms.service", "Missing Topic")
	}

	//Get the records
	// recs, err := h.store.Read("", store.ReadWhere(&store.Fields{
	// "topic": topic,
	// }))
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return nil, errors.InternalServerError("dev.viqchat.rooms.service", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("dev.viqchat.rooms.service", "Room not found")
	}

	// Decode the room
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, r := range recs {
		var room *rooms.Room
		if err := json.Unmarshal(r.Value, &room); err != nil {
			return nil, errors.InternalServerError("dev.viqchat.rooms.service", "Could not unmarshal room: %v", err)
		}
		if room.Topic == topic {
			return room, nil
		}
	}

	return nil, errors.NotFound("dev.viqchat.rooms.service", "Room not found")
}

// List is used to list all rooms for specific user (request can be extended for different params)
func (h *Rooms) List(ctx context.Context, req *rooms.ListRequest, rsp *rooms.Response) error {
	keys, err := h.store.Read(req.UserId, store.ReadPrefix(), store.ReadFrom("viqchat", "user_room"))
	if err != nil {
		return errors.InternalServerError("dev.viqchat.rooms.service", "Could not read from store: %v", err)
	}

	keysSet := make(map[string]struct{})
	for _, v := range keys {
		newKey := strings.TrimPrefix(v.Key, req.UserId)
		keysSet[newKey[1:]] = struct{}{}
	}

	var recs []*store.Record
	for k := range keysSet {
		records, err := h.store.Read(k)
		if err != nil {
			return errors.InternalServerError("dev.viqchat.rooms.service", "Could not read from store: %v", err)
		}

		recs = append(recs, records...)
	}

	// Decode the records
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	rooms := make([]*rooms.Room, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &rooms[i]); err != nil {
			return errors.InternalServerError("dev.viqchat.rooms.service", "Could not unmarshal room: %v", err)
		}
	}

	rsp.Rooms = rooms
	return nil
}
