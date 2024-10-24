package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func DropResetPasswordTokensTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS reset_pw_tokens;"
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop reset_pw_tokens table: %w", err)
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

func DropDocumentsTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS documents;"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop documents table: %w", err)
	}

	return nil
}

func DropSubjectsTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS subjects;"
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop subjects table: %w", err)
	}
	return nil
}

func DropClassesTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS classes;"
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop classes table: %w", err)
	}
	return nil
}

func CreateUsersTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(50) NOT NULL UNIQUE,
        username VARCHAR(20) NOT NULL UNIQUE,
        fullName VARCHAR(50) NOT NULL,
        password VARCHAR(255) NOT NULL,
        gender ENUM('male', 'female') NOT NULL DEFAULT 'male', 
        avatar VARCHAR(255) NOT NULL,
        dateOfBirth DATE NOT NULL,
        role ENUM('user', 'admin') NOT NULL DEFAULT 'user',
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP NULL DEFAULT NULL
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

func CreateResetPasswordTokensTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS reset_pw_tokens (
		token VARCHAR(64) PRIMARY KEY,
		userId INT NOT NULL,
		expiry TIMESTAMP NOT NULL,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create reset_pw_tokens table: %w", err)
	}

	return nil
}

func CreateClassesTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS classes (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create classes table: %w", err)
	}
	return nil
}

func CreateSubjectsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS subjects (
        id INT AUTO_INCREMENT PRIMARY KEY,
        classId INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (classId) REFERENCES classes(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create subjects table: %w", err)
	}
	return nil
}

func CreateDocumentsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS documents (
        id INT AUTO_INCREMENT PRIMARY KEY,
        subjectId INT NOT NULL,
        title VARCHAR(255) NOT NULL,
        fileUrl VARCHAR(255),
        documentType ENUM('PDF', 'VIDEO', 'DOC') NOT NULL DEFAULT 'PDF',
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (subjectId) REFERENCES subjects(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
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
	('user25@gmail.com', 'user25', 'User25', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user31@gmail.com', 'user31', 'User31', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user32@gmail.com', 'user32', 'User32', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user33@gmail.com', 'user33', 'User33', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user34@gmail.com', 'user34', 'User34', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user35@gmail.com', 'user35', 'User35', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user41@gmail.com', 'user41', 'User41', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user42@gmail.com', 'user42', 'User42', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user43@gmail.com', 'user43', 'User43', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user44@gmail.com', 'user44', 'User44', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user'),
	('user45@gmail.com', 'user45', 'User45', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'female', '', '1999-09-09', 'user')
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert test accounts: %w", err)
	}

	return nil
}

func ResetDataBase(db *sql.DB) {
	DropResetPasswordTokensTable(db)
	DropUsersTable(db)
	DropDocumentsTable(db)
	DropSubjectsTable(db)
	DropClassesTable(db)

	CreateUsersTable(db)
	CreateResetPasswordTokensTable(db)
	CreateClassesTable(db)
	CreateSubjectsTable(db)
	CreateDocumentsTable(db)

	InsertTestAccounts(db)
}
