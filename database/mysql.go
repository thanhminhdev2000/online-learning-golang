package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatalf("MYSQL_DSN not set in environment")
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

	// DropUsersTable(db)
	// DropPasswordResetTokensTable(db)
	// err = CreateUsersTable(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create users table: %w", err)
	// }
	// err = CreatePasswordResetTokensTable(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create reset_pw_tokens table: %w", err)
	// }

	return db, nil

}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(50) NOT NULL UNIQUE,
		username VARCHAR(50) NOT NULL UNIQUE,
		fullName VARCHAR(50) NOT NULL,
		password VARCHAR(255) NOT NULL,
		gender ENUM('male', 'female'), 
		avatar VARCHAR(255),
		dateOfBirth DATE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	fmt.Println("Table `users` created successfully!")
	return nil
}

func CreatePasswordResetTokensTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS reset_pw_tokens (
		token VARCHAR(64) PRIMARY KEY,
		userId INT NOT NULL,
		expiry TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create reset_pw_tokens table: %w", err)
	}

	fmt.Println("Table `reset_pw_tokens` created successfully!")
	return nil
}

func DropUsersTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS users;"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop users table: %w", err)
	}

	fmt.Println("Table `users` deleted successfully!")
	return nil
}

func DropPasswordResetTokensTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS reset_pw_tokens;"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop reset_pw_tokens table: %w", err)
	}

	fmt.Println("Table `reset_pw_tokens` deleted successfully!")
	return nil
}
