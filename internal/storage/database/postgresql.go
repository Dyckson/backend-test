package postgres

import (
	config "backend-test/internal/cmd/server"
	"context"
	"log"

	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

func GetDB() *ksql.DB {
	dbConnect, err := kpgx.New(context.Background(), config.DATABASE_URL, ksql.Config{})
	if err != nil {
		log.Panic(err)
	}
	dbConnect.Exec(context.Background(), "set enable_seqscan = off;")

	return &dbConnect
}
