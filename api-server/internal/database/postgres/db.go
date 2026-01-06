package postgres

import (
	"context"
	"log"
	"sync"

	"github.com/ahsansaif47/blockchain-address-watcher/api-server/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Connection *pgx.Conn
	Pool       *pgxpool.Pool
}

var dbInstance *Database
var dbOnce sync.Once

func GetDatabaseInstance() *Database {
	dbOnce.Do(func() {
		dbInstance = Connect()
	})
	return dbInstance
}

func Connect() *Database {
	c := config.GetConfig()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, c.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbPool, err := pgxpool.New(ctx, c.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return &Database{
		Connection: conn,
		Pool:       dbPool,
	}
}
