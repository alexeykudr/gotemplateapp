package backend

type User struct {
	Name     string
	Age      int
	Username string
	Email    string
	IsStuff  bool
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
