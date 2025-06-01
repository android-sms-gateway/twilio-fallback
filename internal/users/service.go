package users

import (
	"errors"

	"github.com/android-sms-gateway/twilio-fallback/internal/encryption"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	GetUser(ID string) (*User, error)
	RegisterUser(login, password, accountSID, authToken string) (*User, error)
	AuthenticateUser(login, password string) (*User, error)
}

type service struct {
	repo      Repository
	encryptor encryption.Encryptor
	logger    *zap.Logger
}

// GetUser implements Service.
func (s *service) GetUser(ID string) (*User, error) {
	return s.repo.GetUser(ID)
}

func NewService(repo Repository, encryptor encryption.Encryptor, logger *zap.Logger) Service {
	return &service{
		repo:      repo,
		encryptor: encryptor,
		logger:    logger,
	}
}

func (s *service) RegisterUser(login, password, accountSID, authToken string) (*User, error) {
	// Check if user already exists
	_, err := s.repo.GetUserBySMSGatewayLogin(login)
	if err == nil {
		return nil, ErrUserAlreadyExists
	} else if !IsUserNotFound(err) {
		s.logger.Error("Error checking if user exists", zap.Error(err))
		return nil, err
	}

	// Encrypt credentials
	encryptedPassword, err := s.encryptor.Encrypt(password)
	if err != nil {
		s.logger.Error("Error encrypting password", zap.Error(err))
		return nil, err
	}

	encryptedAuthToken, err := s.encryptor.Encrypt(authToken)
	if err != nil {
		s.logger.Error("Error encrypting auth token", zap.Error(err))
		return nil, err
	}

	// Generate callback UUID
	callbackUUID := gonanoid.Must(36)

	// Create new user
	user := &User{
		Login:            login,
		Password:         encryptedPassword,
		TwilioAccountSID: accountSID,
		TwilioAuthToken:  encryptedAuthToken,
		CallbackUUID:     callbackUUID,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *service) AuthenticateUser(login, password string) (*User, error) {
	user, err := s.repo.GetUserBySMSGatewayLogin(login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Decrypt and compare password
	decryptedPassword, err := s.encryptor.Decrypt(user.Password)
	if err != nil {
		s.logger.Error("Error decrypting password", zap.Error(err))
		return nil, err
	}

	if decryptedPassword != password {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

var _ Service = (*service)(nil)
