package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDb() (*pgxpool.Pool, error) {
	var db *pgxpool.Pool
	var err error

	for i := range 5 {
		db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err == nil {
			return db, err
		}
		log.Printf("attemt %d/5 Unable to connect Postgres, retry... %v", i+1, err)
		time.Sleep(3 + time.Second)
	}
	return db, err
}
