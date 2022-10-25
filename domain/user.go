package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID       int    `json:"-" db:"id"`
	Email    string `json:"email" binding:"required""`
	Password string `json:"Password" binding:"required"`
}

//Validate user struct to some basic requirements
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Email, validation.Required, validation.Length(5, 50),
			is.Email),
	)
}
