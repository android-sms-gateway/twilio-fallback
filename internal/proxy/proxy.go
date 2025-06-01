package proxy

import "fmt"

type Service interface {
	GetCallbackURL(callbackID string) string
}

type service struct {
	baseURL string
}

func NewService(config Config) Service {
	return &service{
		baseURL: config.BaseURL,
	}
}

func (s *service) GetCallbackURL(callbackID string) string {
	return fmt.Sprintf("%s/callback/%s", s.baseURL, callbackID)
}
