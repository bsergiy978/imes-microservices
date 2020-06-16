package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	vjwt "github.com/COVIDEV/viq-chat-services/jwt/proto/jwt"
	"github.com/dgrijalva/jwt-go"
)

const Secret = "9884c0c0-5a98-401c-ad0f-591f7001278f"

type Jwt struct{}

func (e *Jwt) Generate(ctx context.Context, req *vjwt.GenerateRequest, rsp *vjwt.GenerateResponse) error {
	log.Info("Generating JWT token for ", req.Phone)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{Issuer: req.Phone, Subject: req.Id})

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return err
	}

	rsp.Token = tokenString
	return nil
}
