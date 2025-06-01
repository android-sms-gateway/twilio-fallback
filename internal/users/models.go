package users

import (
	"github.com/android-sms-gateway/twilio-fallback/pkg/core/orm"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Login            string `gorm:"unique;not null;size:16"`
	Password         string `gorm:"not null;size:255"`
	TwilioAccountSID string `gorm:"not null;size:34"`
	TwilioAuthToken  string `gorm:"not null;size:255"`
	CallbackUUID     string `gorm:"unique;not null;size:36"`
}

func init() {
	orm.RegisterMigration(func(db *gorm.DB) error {
		return db.AutoMigrate(&User{})
	})
}
