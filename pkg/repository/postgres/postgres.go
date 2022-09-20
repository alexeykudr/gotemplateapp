package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
	MinConns int32 `default:"20"`
	MaxConns int32 `default:"20"`
	TimeOut  int   `default:"5"`
}

func NewPostgresDB(config PostgresConfig) (*pgxpool.Pool, string, error) {
	ctx, _ := context.WithCancel(context.Background())
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=%d",
		"postgres",
		url.QueryEscape(config.User),
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
		config.TimeOut)

	poolConfig, _ := pgxpool.ParseConfig(connStr)
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConns = config.MaxConns
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Panic(err)
		return nil, connStr, nil
	}

	for i := 0; i < int(config.MaxConns); i++ {
		go func(count int) {
			_, err = pool.Exec(ctx, ";")
			if err != nil {
				log.Error("ping failed", count, err)
			}
		}(i)
		//fmt.Printf("Conections - Max: %d, Iddle: %d, Total: %d \n",
		//	pool.Stat().MaxConns(), pool.Stat().IdleConns(),
		//	pool.Stat().TotalConns())
	}
	return pool, connStr, nil
	//select {}
}
func HealthCheck(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), ";")
	if err != nil {
		return err
	}
	return nil
}
