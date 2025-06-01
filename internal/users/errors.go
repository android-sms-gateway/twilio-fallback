package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

func IsUserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}
