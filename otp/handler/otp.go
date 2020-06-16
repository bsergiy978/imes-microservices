// Package otp provides service for generating OTP passcode.
package otp

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/cockroach"
	"github.com/pquerna/otp/totp"

	otp "github.com/COVIDEV/viq-chat-services/otp/proto/otp"
)

const (
	secretBytesCount = 20
)

// Otp Service.
type Otp struct {
	store store.Store
}

// NewOtpHandler creating OTP service handler.
func NewOtpHandler(srv micro.Service) *Otp {
	return &Otp{store: cockroach.NewStore(store.Database("viqchat"), store.Table("otp_secrets"))}
}

// Generate is generating OTP passcode.
func (e *Otp) Generate(ctx context.Context, req *otp.GenerateRequest, rsp *otp.GenerateResponse) error {
	log.Info("Generatinig passcode for ", req.Phone)
	secret := make([]byte, secretBytesCount)
	_, err := rand.Read(secret)
	if err != nil {
		return err
	}

	passcode, err := totp.GenerateCode(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), time.Now())
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info(e.store.Options().Table)
	err = e.store.Write(&store.Record{Key: req.Phone, Value: secret})
	if err != nil {
		log.Error(err.Error())
		return err
	}

	rsp.Passcode = passcode
	return nil
}

// Verify is verifying OTP passcode.
func (e *Otp) Verify(ctx context.Context, req *otp.VerifyRequest, rsp *otp.VerifyResponse) error {
	log.Info("Verifying passcode from ", req.Phone)
	recs, err := e.store.Read(req.Phone)
	if err != nil {
		return err
	}

	for _, r := range recs {
		if r.Key != req.Phone {
			continue
		}

		if valid := totp.Validate(req.Passcode, base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(r.Value)); valid {
			log.Info(req.Phone, " verified successfully")
		} else {
			log.Info("Failed to verify ", req.Phone)
			return errors.BadRequest("com.github.romatroskin.viqchat.otp.service.otp", "Failed to verify passcode")
		}
	}

	return nil
}
