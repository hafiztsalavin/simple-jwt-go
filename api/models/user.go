package models

import (
	"errors"
	"simple-jwt-go/api/utils"
	"time"

	validator "github.com/asaskevich/govalidator"
)

type User struct {
	ID        uint32    `gorm:"serial;" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rr *RegistrationRequest) CreateUser(target *User) error {
	if !utils.IsNonEmpty(rr.Username, rr.Email, rr.Password) {
		return errors.New("empty format")
	}

	if !validator.IsEmail(rr.Email) || !validator.IsPrintableASCII(rr.Username) || !validator.IsPrintableASCII(rr.Password) {
		return errors.New("bad format character")
	}

	target.Email = rr.Email
	target.Username = rr.Username
	target.Password = rr.Password

	return nil
}
