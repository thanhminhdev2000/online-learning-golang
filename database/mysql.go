package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL() (*sql.DB, error) {
	dsn := os.Getenv("DB_CONNECTION")
	if dsn == "" {
		log.Fatalf("DB_CONNECTION not set in environment")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open MySQL connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not ping MySQL: %w", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	return db, nil

}
