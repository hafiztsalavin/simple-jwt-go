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

type UserResponse struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"registered_at"`
}

type UserUpdateRequest struct {
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

func (ur *UserResponse) InsertFromModel(user User) error {
	ur.ID = user.ID
	ur.Username = user.Username
	ur.Email = user.Email
	ur.CreatedAt = user.CreatedAt
	return nil
}

func (ur *UserUpdateRequest) UpdateUserModel(target *User) error {
	if !validator.IsPrintableASCII(ur.Username) || !validator.IsPrintableASCII(ur.Password) || !validator.IsEmail(ur.Email) {
		return errors.New("email/password have bad format or character")
	}

	if ur.Email != "" {
		target.Email = ur.Email
	}

	if ur.Password != "" {
		target.Password = ur.Password
	}

	if ur.Username != "" {
		target.Username = ur.Username
	}

	return nil
}
