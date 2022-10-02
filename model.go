package backend

import "time"

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
	Email    string `json:"email" binding:"required""`
	IsStuff  bool   `json:"isStuff" binding:"required"`
}

type Proxy struct {
	Dhcp int
}

type Client struct {
	ContactLogin string
	EndDate      time.Time
	OrderedProxy []Proxy
	Payment      float32
}

type Ports struct {
	OrderId    []int `json:"orderId" bson:"orderId"`
	ReservedAt int64 `json:"reservedAt" bson:"reservedAt"`
}

type MongoUser struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Ports
}
