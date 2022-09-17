package main

import (
	"backend/internal/service"
	"backend/pkg/repository"
	"fmt"
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
	fmt.Println(viper.GetString("TOKEN"))

	pool, err := repository.NewPostgresDB(repository.Config{
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

	err = repository.Apply(pool)
	if err != nil {
		fmt.Println(err)
	}

	ins := service.NewInstance(pool)
	ins.Start()

}
