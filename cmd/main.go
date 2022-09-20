package main

import (
	"backend/pkg/repository/mongo"
	"backend/pkg/repository/postgres"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := NewConfig()
	if err != nil {
		log.Panic(err)
		panic("Error with config!")
	}
	log.SetFormatter(&log.JSONFormatter{})

	pool, ConnString, err := postgres.NewPostgresDB(postgres.PostgresConfig{
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
		MinConns: 10,
		MaxConns: 20,
		TimeOut:  5,
	},
	)

	err = postgres.HealthCheck(pool)
	if err != nil {
		fmt.Println(err)
	}

	ins := postgres.NewInstance(pool)
	ins.Start()

	mdb, err := sql.Open("postgres", ConnString)
	if err != nil {
		fmt.Println(err)
	}
	err = mdb.Ping()
	if err != nil {
		panic(err)
	}
	//err = goose.Up(mdb, "internal/migrations")
	//err = goose.Down(mdb, "internal/migrations")
	//if err != nil {
	//	panic(err)
	//}
	//
	mongoClient, err := mongo.NewMongoDbClient(mongo.MongoConfig{
		User:     viper.GetString("MONGO_USER"),
		Password: viper.GetString("MONGO_PASSWORD"),
		Host:     viper.GetString("MONGO_HOST"),
		Port:     viper.GetInt("MONGO_PORT"),
	})

	err = mongo.HealthCheck(mongoClient)

	mongoInstance := mongo.NewInstance(mongoClient)
	mongoInstance.Start()

}
