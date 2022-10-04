package backend

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
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
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&u.Password, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&u.Email, validation.Required, validation.Length(5, 50)),
	)
}

//Proxy model of proxy instance in DB
type Proxy struct {
	Dhcp      int
	TelNumber int
}

//Client Model of client instance in DB
type Client struct {
	ContactLogin string
	EndDate      time.Time
	OrderedProxy []Proxy
	Payment      float32
}

//model of ports instance in DB
type Ports struct {
	OrderId    []int `json:"orderId" bson:"orderId"`
	ReservedAt int64 `json:"reservedAt" bson:"reservedAt"`
}

//Tmp mongo instance
type MongoUser struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Ports
}
