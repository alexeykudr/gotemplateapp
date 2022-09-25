package mongo

import (
	"backend"
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Instance struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewInstance(client *mongo.Client) *Instance {
	return &Instance{client: client}
}

func (i *Instance) Start() {
	//i.FindByUsername("clients", "proxy_clients", "petya")
	err := i.FindByName("user2")
	if err != nil {
		return
	}
}
func (i *Instance) MockSingeDocument() error {
	mongoUser1 := backend.MongoUser{
		Username: "petya",
		Email:    "petia123123123@gmail.com",
		Ports: backend.Ports{
			OrderId:    []int{1, 2, 3, 4, 11, 12, 19},
			ReservedAt: time.Now().Unix(),
		},
	}
	var batch []interface{}
	batch = append(batch, mongoUser1)

	collection, err := i.InsertManyToCollection("clients", "proxy_clients", batch)
	if err != nil {
		return err
	}
	fmt.Println(collection)
	return nil
}

func (i *Instance) InsertManyToCollection(db, collection string, batch []interface{}) ([]interface{}, error) {
	cl := i.client.Database(db).Collection(collection)

	insertResult, err := cl.InsertMany(context.TODO(), batch)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedIDs)
	return insertResult.InsertedIDs, nil
}
func (i *Instance) FindByUsername(db, collection, username string) error {
	filter := bson.M{
		"username": username,
	}
	var mongoUser backend.MongoUser
	err := i.client.Database(db).Collection(collection).FindOne(context.Background(), filter).Decode(&mongoUser)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(mongoUser)
	return nil
}

func (i *Instance) FindByName(n string) error {
	filter := bson.M{"name": "user2"}
	findOptions := options.Find()
	var users []interface{}

	rows, err := i.client.Database("clients").Collection("proxy_clients").Find(context.TODO(), filter, findOptions)
	if err != nil {
		return err
	}
	for rows.Next(context.TODO()) {
		var user interface{}
		err := rows.Decode(&user)
		fmt.Println(user)
		if err != nil {
			return err
		}
		users = append(users, &user)
	}
	fmt.Println(users)
	rows.Close(context.Background())
	return nil
}

//TODO https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/
//TODO Переделать на эирбнб , реализовать круд , пару агреграции, поиск
//GET BY ID, GET BY MIN NIGHTS TO MAX NIGHTS , GET BY accommodates (вмещает) , GET BY PRICE
