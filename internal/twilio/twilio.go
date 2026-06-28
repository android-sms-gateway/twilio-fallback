package twilio

import (
	"context"
	"fmt"

	"github.com/android-sms-gateway/twilio-fallback/internal/common"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Service interface {
	GetMessage(ctx context.Context, sid string) (common.Message, error)
	ValidateSignature(url string, params map[string]string, signature string) error
}

type service struct {
	accountSID  string
	callbackURL string

	client    *twilio.RestClient
	validator client.RequestValidator
}

func NewService(config Config) Service {
	return &service{
		accountSID:  config.AccountSID,
		callbackURL: config.CallbackURL,

		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:                 config.AccountSID,
			Password:                 config.AuthToken,
			AccountSid:               "",
			Client:                   nil,
			ClientCredentialProvider: nil,
		}),
		validator: client.NewRequestValidator(config.AuthToken),
	}
}

func (s *service) GetMessage(_ context.Context, sid string) (common.Message, error) {
	var params openapi.FetchMessageParams
	resp, err := s.client.Api.FetchMessage(sid, &params)
	if err != nil {
		return common.Message{}, fmt.Errorf("twilio fetch message: %w", err)
	}

	if resp.To == nil || resp.Body == nil {
		return common.Message{}, ErrMissingRequiredFields
	}

	return common.Message{
		ID:   *resp.Sid,
		To:   *resp.To,
		Body: *resp.Body,
	}, nil
}

func (s *service) ValidateSignature(url string, params map[string]string, signature string) error {
	// Validate AccountSID
	accountSid, ok := params["AccountSid"]
	if !ok {
		return ErrMissingAccountSid
	}
	if accountSid != s.accountSID {
		return ErrAccountSidMismatch
	}

	if s.callbackURL != "" {
		url = s.callbackURL
	}

	// Validate signature
	if !s.validator.Validate(url, params, signature) {
		return ErrSignatureValidationFailed
	}

	return nil
}

var _ Service = (*service)(nil)
