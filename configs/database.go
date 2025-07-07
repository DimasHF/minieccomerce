package configs

import (
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

var once sync.Once

func InitDB() *sqlx.DB {
	once.Do(func() {
		host := "127.0.0.1"
		port := "5432"
		user := "minieccomerce"
		psw := "minieccomerce"
		db := "minieccomercedb"
		urldb := "postgres" + "://" + user + ":" + psw + "@" + host + ":" + port + "/" + db

		conn, err := sqlx.Connect("pgx", urldb)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		DB = conn
		DB.SetMaxOpenConns(100)
		DB.SetMaxIdleConns(10)
		DB.SetConnMaxLifetime(time.Hour)
	})

	return DB
}
