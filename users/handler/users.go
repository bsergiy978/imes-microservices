package users

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/segmentio/ksuid"

	pb "github.com/COVIDEV/viq-chat-services/users/proto/users"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/cockroach"
)

var (
	// URLSafeRegex is a function which returns true if a string is URL safe
	URLSafeRegex = regexp.MustCompile(`^[A-Za-z0-9_-].*?$`).MatchString
)

// Handler implements the users service interface
type Handler struct {
	store     store.Store
	publisher micro.Publisher
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	// Return the initialised store
	return &Handler{
		store: cockroach.NewStore(
			store.Database("viqchat"),
			store.Table("users"),
		),
		publisher: micro.NewPublisher(srv.Name(), srv.Client()),
	}
}

// Create inserts a new user into the store
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Validate request
	if req.User == nil || req.User.Phone == "" {
		return errors.BadRequest("dev.viqchat.users.service", "User is missing")
	}

	// Check to see if the user already exists
	if usr, err := h.findUserByPhone(req.User.Phone); err == nil {
		//TODO: proper error handling
		// return errors.Conflict("dev.viqchat.users.service", "User already exist")
		rsp.User = usr
		return nil
	} else {
		log.Error(err.Error())
	}

	// If validation only, return here
	if req.ValidateOnly {
		return nil
	}

	// Add the auto-generated fields
	var user pb.User = *req.User
	if len(user.Id) == 0 {
		// generate ID
		user.Id = ksuid.New().String()
	}

	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	// Encode the user
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes, Metadata: cockroach.Metadata{
		"phone": user.Phone,
	}}); err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not write to store: %v", err)
	}

	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserCreated,
	})

	// Return to the user and token in the response
	rsp.User = &user
	return nil
}

// Read retrieves a user from the store
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	var user *pb.User
	var err error

	if len(req.Phone) > 0 {
		user, err = h.findUserByPhone(req.Phone)
	} else {
		user, err = h.findUser(req.Id)
	}

	if err != nil {
		return err
	}

	rsp.User = user
	return nil
}

// Update modifies a user in the store
func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Validate the request
	if req.User == nil {
		return errors.BadRequest("dev.viqchat.users.service", "User is missing")
	}

	// Lookup the user
	user, err := h.findUser(req.User.Id)
	if err != nil {
		return err
	}

	// Update the user with the given attributes
	// TODO: Find a way which allows only updating a subset of attributes,
	// checking for blank values doesn't work since there needs to be a way
	// of unsetting attributes.
	user.Username = req.User.Username
	user.Displayname = req.User.Displayname
	user.Bio = req.User.Bio
	user.Updated = time.Now().Unix()

	// Encode updated user
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes}); err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserUpdated,
		User: user,
	})

	// Return the user in the response
	rsp.User = user
	return nil
}

// Delete a user in the store
func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Lookup the user
	user, err := h.findUser(req.Id)
	if err != nil {
		return err
	}

	// Delete from the store
	if err := h.store.Delete(user.Id); err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserDeleted,
		User: user,
	})

	return nil
}

// Search the users in th store, using full name
func (h *Handler) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	// List all the records
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("dev.viqchat.users.service", "Could not read from store: %v", err)
	}

	// Decode the records
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	users := make([]*pb.User, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &users[i]); err != nil {
			return errors.InternalServerError("dev.viqchat.users.service", "Could not unmarshal user: %v", err)
		}
	}

	// Filter and return the users
	rsp.Users = make([]*pb.User, 0)
	for _, u := range users {
		anyname := fmt.Sprintf("%v %v", u.Displayname, u.Username)
		if strings.Contains(anyname, req.Query) {
			rsp.Users = append(rsp.Users, u)
		}
	}

	return nil
}

// findUser retreives a user given an ID. It is used by the Read, Update
// and Delete functions
func (h *Handler) findUser(id string) (*pb.User, error) {
	// Validate the request
	if len(id) == 0 {
		return nil, errors.BadRequest("dev.viqchat.users.service", "Missing ID")
	}

	// Get the records
	recs, err := h.store.Read(id)
	if err != nil {
		return nil, errors.InternalServerError("dev.viqchat.users.service", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("dev.viqchat.users.service", "User not found")
	}
	if len(recs) > 1 {
		return nil, errors.InternalServerError("dev.viqchat.users.service", "Store corrupted, %b records found for ID", len(recs))
	}

	// Decode the user
	var user *pb.User
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return nil, errors.InternalServerError("dev.viqchat.users.service", "Could not unmarshal user: %v", err)
	}

	return user, nil
}

// findUserByPhone retreives a user given a phone
func (h *Handler) findUserByPhone(phone string) (*pb.User, error) {
	// Validate request
	if len(phone) == 0 {
		return nil, errors.BadRequest("dev.viqchat.users.service", "Missing Phone")
	}

	//Get the records
	// recs, err := h.store.Read("", store.ReadWhere(&store.Fields{
	// "phone": phone,
	// }))
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return nil, errors.InternalServerError("dev.viqchat.users.service", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("dev.viqchat.users.service", "Users not found")
	}

	// Decode the user
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, r := range recs {
		var user *pb.User
		if err := json.Unmarshal(r.Value, &user); err != nil {
			return nil, errors.InternalServerError("dev.viqchat.users.service", "Could not unmarshal user: %v", err)
		}
		if user.Phone == phone {
			return user, nil
		}
	}

	return nil, errors.NotFound("dev.viqchat.users.service", "User not found")
}
