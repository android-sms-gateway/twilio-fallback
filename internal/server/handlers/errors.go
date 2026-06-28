package handlers

import "errors"

var ErrProxyServiceNil = errors.New("proxy service is nil")
var ErrTwilioServiceNil = errors.New("twilio service is nil")
var ErrValidatorNil = errors.New("validator is nil")
var ErrLoggerNil = errors.New("logger is nil")
