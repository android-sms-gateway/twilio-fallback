package proxy

import (
	"context"
	"fmt"

	"github.com/android-sms-gateway/twilio-fallback/internal/smsgate"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"go.uber.org/zap"
)

type Service interface {
	Start()
	ProcessCallback(ctx context.Context, messageSid, messageStatus string) error
	Stop()
}

type service struct {
	twilio  twilio.Service
	smsGate smsgate.Service

	logger *zap.Logger

	jobs *jobsService
}

func NewService(twilio twilio.Service, smsGate smsgate.Service, logger *zap.Logger) Service {
	return &service{
		twilio:  twilio,
		smsGate: smsGate,
		logger:  logger,

		jobs: newJobsService(logger),
	}
}

func (s *service) Start() {
	s.jobs.Start()
}

func (s *service) ProcessCallback(ctx context.Context, messageSid, messageStatus string) error {
	// Filter non-failure statuses
	if !isFailedStatus(messageStatus) {
		return nil
	}

	// Enqueue job with retry operation
	err := s.jobs.Enqueue(messageSid, func(ctx context.Context) error {
		// Fetch full message from Twilio
		msg, err := s.twilio.GetMessage(ctx, messageSid)
		if err != nil {
			return fmt.Errorf("failed to fetch message from Twilio: %w", err)
		}

		// Retry with exponential backoff
		_, err = s.smsGate.Send(ctx, msg)
		if err != nil {
			return fmt.Errorf("failed to send message to SMS gate: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

func isFailedStatus(status string) bool {
	return status == "failed" || status == "undelivered"
}

func (s *service) Stop() {
	s.jobs.Close()
}

var _ Service = (*service)(nil)
