package main

import (
	"backend"
	"backend/pkg/handler"
	"backend/pkg/repository/postgres"
	"backend/pkg/service"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func NewConfig() error {
	log.SetFormatter(&log.JSONFormatter{})
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func BasicAuthHandler(w http.ResponseWriter, r *http.Request) {
	a := "qwe"
	b := "123"
	username, password, ok := r.BasicAuth()
	if ok {
		if username == a && password == b {
			fmt.Fprint(w, "ok!")
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
func ReturnIdsHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := r.URL.Query()["id[]"]
	if !ok {
		fmt.Fprint(w, "empty id")
	} else {
		user := backend.User{
			Name:     "user3",
			Age:      11,
			Username: "user3",
			Email:    "user3@gmail.com",
			IsStuff:  false,
		}
		jsondata, _ := json.Marshal(user)
		_, _ = w.Write(jsondata)

	}

}

func main() {
	err := NewConfig()
	if err != nil {
		log.Panic(err)
		panic("Error with config!")
	}

	pool, _, err := postgres.NewPostgresDB(postgres.PostgresConfig{
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

	repo := postgres.NewRepository(pool)
	services := service.NewService(repo)
	handler := handler.NewHandler(services)

	//mongoClient, err := mongo.NewMongoDbClient(mongo.Config{
	//	User:     viper.GetString("MONGO_USER"),
	//	Password: viper.GetString("MONGO_PASSWORD"),
	//	Host:     viper.GetString("MONGO_HOST"),
	//	Port:     viper.GetInt("MONGO_PORT"),
	//})
	//fmt.Println(mongoClient)

	//m := mongo.NewAirbnbMongoInstance(mongoClient)
	//m.FindByType("10038496")

	//http.HandleFunc("/", BasicAuthHandler)
	//http.HandleFunc("/orders", ReturnIdsHandler)
	http.ListenAndServe(":8080", handler.InitRoutes())
}
