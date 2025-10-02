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

// GetDB retorna uma instância singleton da conexão com o banco
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

// CloseDB fecha a conexão com o banco (usar apenas no shutdown da aplicação)
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
