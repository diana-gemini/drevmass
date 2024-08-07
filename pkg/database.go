package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConn() (*pgxpool.Pool, error) {
	connectionString := os.Getenv("DNS")
	if connectionString == "" {
		return nil, fmt.Errorf("DNS environment variable not set")
	}

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	fmt.Println("Success connection")

	return pool, nil
}
