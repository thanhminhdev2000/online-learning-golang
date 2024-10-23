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

	// DropPasswordResetTokensTable(db)
	// DropUsersTable(db)
	// CreateUsersTable(db)
	// CreatePasswordResetTokensTable(db)
	// InsertTestAccounts(db)

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
		gender ENUM('male', 'female') NOT NULL DEFAULT 'male', 
		avatar VARCHAR(255) NOT NULL,
		dateOfBirth DATE NOT NULL,
		role ENUM('user', 'admin') NOT NULL DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

func InsertTestAccounts(db *sql.DB) error {
	query := `
	INSERT INTO users (email, username, fullName, password, gender, avatar, dateOfBirth, role)
	VALUES ('admin1@gmail.com', 'admin1', 'Admin1', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '', '1985-01-01', 'admin'),
	('admin2@gmail.com', 'admin2', 'Admin2', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '', '1985-01-01', 'admin'),
	('admin3@gmail.com', 'admin3', 'Admin3', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '', '1985-01-01', 'admin'),
	('admin4@gmail.com', 'admin4', 'Admin4', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '', '1985-01-01', 'admin'),
	('admin5@gmail.com', 'admin5', 'Admin5', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '', '1985-01-01', 'admin'),
	('user01@gmail.com', 'user01', 'User01', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user02@gmail.com', 'user02', 'User02', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user03@gmail.com', 'user03', 'User03', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user04@gmail.com', 'user04', 'User04', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user05@gmail.com', 'user05', 'User05', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user11@gmail.com', 'user11', 'User11', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user12@gmail.com', 'user12', 'User12', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user13@gmail.com', 'user13', 'User13', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user14@gmail.com', 'user14', 'User14', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user15@gmail.com', 'user15', 'User15', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user21@gmail.com', 'user21', 'User21', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user22@gmail.com', 'user22', 'User22', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user23@gmail.com', 'user23', 'User23', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user24@gmail.com', 'user24', 'User24', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user25@gmail.com', 'user25', 'User25', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user')
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert test accounts: %w", err)
	}

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

	return nil
}

func DropUsersTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS users;"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop users table: %w", err)
	}

	return nil
}

func DropPasswordResetTokensTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS reset_pw_tokens;"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop reset_pw_tokens table: %w", err)
	}

	return nil
}
