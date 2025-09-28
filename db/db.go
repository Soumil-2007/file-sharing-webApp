package db

import (
	"database/sql"
	"log"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

var DB *sql.DB

func Init() error {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")   // e.g. 127.0.0.1
	port := os.Getenv("DB_PORT")   // e.g. 3306
	name := os.Getenv("DB_NAME")   // set to filesharing

	if user == "" || name == "" {
		return fmt.Errorf("DB_USER and DB_NAME environment variables are required")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	return DB.Ping()
}