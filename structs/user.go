package structs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID       int    `json:"-" db:"id"`
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
	Email    string `json:"email" binding:"required""`
	IsStuff  bool   `json:"isStuff" binding:"required"`
}

//Validate user struct to some basic requirements
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Email, validation.Required, validation.Length(5, 50)),
	)
}
