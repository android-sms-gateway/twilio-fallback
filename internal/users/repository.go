package users

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetUser(ID string) (*User, error)
	CreateUser(user *User) error
	GetUserBySMSGatewayLogin(login string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUser(ID string) (*User, error) {
	var user User
	err := r.db.Where("login = ?", ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetUserBySMSGatewayLogin(login string) (*User, error) {
	var user User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

var _ Repository = (*repository)(nil)
