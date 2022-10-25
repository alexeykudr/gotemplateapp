package domain

import "time"

//Client Model of client instance in DB
type Client struct {
	ContactLogin string
	EndDate      time.Time
	OrderedProxy []Proxy
	Payment      float32
}
