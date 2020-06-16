package handler

import (
	"net/http"
	"time"

	jwt "github.com/COVIDEV/viq-chat-services/jwt/proto/jwt"
	otp "github.com/COVIDEV/viq-chat-services/otp/proto/otp"
	users "github.com/COVIDEV/viq-chat-services/users/proto/users"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/v2/store"
)

// Auth is used to request the code to phone number
type Auth struct {
	Phone string `from:"phone" json:"phone" xml:"phone" binding:"required"`
}

// Validate is used to validate the code sent to phone number
type Validate struct {
	Phone      string `from:"phone" json:"phone" xml:"phone" binding:"required"`
	Code       string `from:"code" json:"code" xml:"code" binding:"required"`
	DeviceID   string `from:"deviceId" json:"deviceId" xml:"deviceId" binding:"required"`
	DeviceName string `from:"deviceName" json:"deviceName" xml:"deviceName" binding:"required"`
}

type Access struct {
	UserID   string `json:"user_id,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
	Created  int64  `json:"created,omitempty"`
	Updated  int64  `json:"updated,omitempty"`
}

// Handler implements the user api interface
type Handler struct {
	Store store.Store
	Otp   otp.OtpService
	Jwt   jwt.JwtService
	Users users.UsersService
}

// Auth is used to request confirmation code
func (h *Handler) Auth(ctx *gin.Context) {
	var json Auth
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Request validation failed!",
			"error":   "Phone number is required for using our app.",
		})
		return
	}

	// userResp, err := h.Users.Read(ctx, &users.ReadRequest{Phone: json.Phone})
	// if userResp != nil {
	// 	ctx.JSON(http.StatusConflict, gin.H{
	// 		"message": "Request validation failed!",
	// 		"error":   "User already exists!",
	// 	})
	// 	return
	// }

	resp, err := h.Otp.Generate(ctx, &otp.GenerateRequest{Phone: json.Phone})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"passcode": resp.Passcode,
	})
}

// Verify is used to verify otp code with otp service
func (h *Handler) Verify(ctx *gin.Context) {
	var json Validate
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errors.BadRequest("com.github.romatroskin.viqchat.auth.api", "Mising required fields."),
		)
		return
	}

	_, err := h.Otp.Verify(ctx, &otp.VerifyRequest{Phone: json.Phone, Passcode: json.Code})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userResp, err := h.Users.Create(ctx, &users.CreateRequest{User: &users.User{Phone: json.Phone}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	resp, err := h.Jwt.Generate(ctx, &jwt.GenerateRequest{Id: userResp.User.Id, Phone: json.Phone, Secret: json.Code})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var access Access = Access{UserID: userResp.User.Id, Created: time.Now().Unix(), Updated: time.Now().Unix()}
	var j = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := j.Marshal(access)
	if err = h.Store.Write(&store.Record{Key: resp.Token, Value: bytes}); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  userResp.GetUser(),
		"token": resp.Token,
	},
	)
}

func (h *Handler) Profile(ctx *gin.Context) {
	recs, err := h.Store.Read(ctx.Request.Header.Get("x-auth-token"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var access *Access
	var j = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := j.Unmarshal(recs[0].Value, &access); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	userResp, err := h.Users.Read(ctx, &users.ReadRequest{Id: access.UserID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userResp.GetUser(),
	},
	)
}

func (h *Handler) Contacts(ctx *gin.Context) {
	resp, err := h.Users.Search(ctx, &users.SearchRequest{Query: ""})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"contacts": resp.GetUsers(),
	},
	)
}
