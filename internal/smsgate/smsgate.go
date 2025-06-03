package smsgate

import (
	"context"
	"fmt"
	"time"

	"github.com/android-sms-gateway/client-go/rest"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/twilio-fallback/internal/common"
	"github.com/cenkalti/backoff/v5"
	"go.uber.org/zap"
)

type Service interface {
	Send(ctx context.Context, msg common.Message) (string, error)
}

type service struct {
	timeout time.Duration

	client *smsgateway.Client

	logger *zap.Logger
}

func NewService(config Config, logger *zap.Logger) (Service, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	return &service{
		timeout: config.Timeout,

		client: smsgateway.NewClient(smsgateway.Config{
			BaseURL:  config.BaseURL,
			User:     config.Username,
			Password: config.Password,
		}),

		logger: logger,
	}, nil
}

func (c *service) Send(ctx context.Context, msg common.Message) (string, error) {
	sms := smsgateway.Message{
		ID:           msg.ID,
		Message:      msg.Body,
		PhoneNumbers: []string{msg.To},
	}

	operation := func() (string, error) {
		ctx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()

		state, err := c.client.Send(ctx, sms)
		if rest.IsConflict(err) {
			c.logger.Info("Message already enqueued", zap.String("id", msg.ID))
			return msg.ID, nil
		}
		if rest.IsClientError(err) {
			return "", backoff.Permanent(fmt.Errorf("smsgate send: %w", err))
		}

		if err != nil {
			return "", fmt.Errorf("smsgate send: %w", err)
		}

		c.logger.Info("Message enqueued", zap.String("id", state.ID), zap.String("state", string(state.State)))
		return state.ID, nil
	}

	result, err := backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return "", fmt.Errorf("smsgate retry: %w", err)
	}

	return result, nil
}

var _ Service = (*service)(nil)
