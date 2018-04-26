package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	DB_URL = os.Getenv("DB_URL")
)
var db *sql.DB
var once sync.Once

func Connect() *sql.DB {
	once.Do(func() {
		//local
		connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			"localhost",
			"apisit",
			"o3_neo_dev",
			"")
		var err error
		if DB_URL == "" {
			DB_URL = connectionString
		}
		db, err = sql.Open("postgres", DB_URL)
		if err != nil {
			log.Printf("cannot open database %v", err)
			db = nil
		}
		db.SetMaxOpenConns(256)
		// db.SetConnMaxLifetime(time.Hour * 24)
		db.SetMaxIdleConns(256)
	})
	return db
}
