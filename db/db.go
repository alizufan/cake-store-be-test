package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Init() *sql.DB {
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		user = "root"
	}

	pass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		pass = "secret"
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		host = "localhost"
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		port = "3306"
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		name = "cake-store"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open db connection: \n%+v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db connection: \n%+v\n", err)
	}

	return db
}
