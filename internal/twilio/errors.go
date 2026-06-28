package twilio

import "errors"

var ErrMissingRequiredFields = errors.New("twilio response missing required fields")
var ErrMissingAccountSid = errors.New("missing AccountSid parameter")
var ErrAccountSidMismatch = errors.New("AccountSid mismatch")
var ErrSignatureValidationFailed = errors.New("twilio signature validation failed")
