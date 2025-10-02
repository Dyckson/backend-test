package postgres

import (
	config "backend-test/internal/cmd/server"
	"context"
	"log"
	"sync"

	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

var (
	db   *ksql.DB
	once sync.Once
)

func GetDB() *ksql.DB {
	once.Do(func() {
		dbConnect, err := kpgx.New(context.Background(), config.DATABASE_URL, ksql.Config{})
		if err != nil {
			log.Panic(err)
		}
		dbConnect.Exec(context.Background(), "set enable_seqscan = off;")
		db = &dbConnect
	})
	return db
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
