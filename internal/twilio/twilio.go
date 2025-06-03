package twilio

import (
	"context"
	"errors"
	"fmt"

	"github.com/android-sms-gateway/twilio-fallback/internal/common"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Service interface {
	GetMessage(ctx context.Context, sid string) (common.Message, error)
}

type service struct {
	client *twilio.RestClient
}

func NewService(config Config) Service {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &service{
		client: client,
	}
}

func (s *service) GetMessage(ctx context.Context, sid string) (common.Message, error) {
	params := &openapi.FetchMessageParams{}
	resp, err := s.client.Api.FetchMessage(sid, params)
	if err != nil {
		return common.Message{}, fmt.Errorf("twilio fetch message: %w", err)
	}

	if resp.To == nil || resp.Body == nil {
		return common.Message{}, errors.New("twilio response missing required fields")
	}

	return common.Message{
		ID:   *resp.Sid,
		To:   *resp.To,
		Body: *resp.Body,
	}, nil
}

var _ Service = (*service)(nil)
