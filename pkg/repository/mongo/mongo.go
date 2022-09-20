package mongo

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func NewMongoConfig(user string, password string, port int) *MongoConfig {
	return &MongoConfig{User: user, Password: password, Port: port}
}

func NewMongoDbClient(config MongoConfig) (*mongo.Client, error) {
	connect := fmt.Sprintf("mongodb://%s:%s@%s:%d/", config.User, config.Password, config.Host, config.Port)
	//mongodb://root:example@localhost:27017/
	client, _ := mongo.NewClient(options.Client().ApplyURI(connect))
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Check the connection

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func HealthCheck(client *mongo.Client) error {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
