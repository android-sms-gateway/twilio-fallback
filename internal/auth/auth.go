package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type AuthService struct {
	secret string
	expiry time.Duration

	logger *zap.Logger
}

func NewAuthService(logger *zap.Logger, secret string, expiry time.Duration) *AuthService {
	if logger == nil {
		panic("logger is nil")
	}

	if secret == "" {
		panic("secret is empty")
	}

	if expiry == 0 {
		panic("expiry is zero")
	}

	return &AuthService{
		logger: logger,
		secret: secret,
		expiry: expiry,
	}
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"jti": gonanoid.Must(),
		"sub": userID,
		"exp": time.Now().Add(s.expiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, exists := claims["sub"]; exists {
			if userID, ok := sub.(string); ok {
				return userID, nil
			}
		}
		return "", ErrInvalidToken
	}

	return "", ErrInvalidToken
}
