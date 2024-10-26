package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
        avatar VARCHAR(255) DEFAULT "",
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
        documentType ENUM('PDF', 'VIDEO') NOT NULL DEFAULT 'PDF',
		views INT DEFAULT 0,
        downloads INT DEFAULT 0,
		author VARCHAR(255) DEFAULT "",
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
	INSERT INTO users (email, username, ?, password, gender, dateOfBirth, role)
	VALUES (?, ?, 'Admin', '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '1985-01-01', ?)
	`

	for index := 1; index < 10; index++ {
		role := "user"
		fullName := "User Name"
		if index%2 == 1 {
			role = "admin"
			fullName = "Admin Name"
		}

		IdStr := strconv.Itoa(index)
		_, err := db.Exec(query, role+IdStr+"@gmail.com", role+IdStr, fullName+" "+IdStr, role)

		if err != nil {
			return fmt.Errorf("failed to insert document")
		}
	}

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert test accounts: %w", err)
	}

	return nil
}

func InsertClassesData(db *sql.DB) error {
	classes := []string{
		"Đề thi thử đại học", "Lớp 12", "Lớp 11", "Lớp 10", "Thi vào lớp 10",
		"Lớp 9", "Lớp 8", "Lớp 7", "Lớp 6", "Thi vào lớp 6", "Lớp 5", "Lớp 4",
		"Lớp 3", "Lớp 2", "Lớp 1",
	}

	query := `INSERT INTO classes (name) VALUES (?)`
	for _, className := range classes {
		_, err := db.Exec(query, className)
		if err != nil {
			return fmt.Errorf("failed to insert class %s: %w", className, err)
		}
	}

	return nil
}

func InsertSubjectsData(db *sql.DB) error {
	subjects := []string{
		"Ngữ Văn", "Toán", "Tiếng Anh", "Vật lí", "Hóa Học", "Sinh học",
		"Lịch sử", "Địa lí", "Giáo dục công dân",
	}

	queryGetClassId := `SELECT id FROM classes WHERE name = ?`
	queryInsertSubject := `INSERT INTO subjects (classId, name) VALUES (?, ?)`

	classes := []string{
		"Đề thi thử đại học", "Lớp 12", "Lớp 11", "Lớp 10", "Thi vào lớp 10",
		"Lớp 9", "Lớp 8", "Lớp 7", "Lớp 6", "Thi vào lớp 6", "Lớp 5", "Lớp 4",
		"Lớp 3", "Lớp 2", "Lớp 1",
	}

	for _, className := range classes {
		var classId int
		err := db.QueryRow(queryGetClassId, className).Scan(&classId)
		if err != nil {
			return fmt.Errorf("failed to retrieve classId for class %s: %w", className, err)
		}

		for _, subject := range subjects {
			_, err := db.Exec(queryInsertSubject, classId, subject)
			if err != nil {
				return fmt.Errorf("failed to insert subject %s for class %s: %w", subject, className, err)
			}
		}
	}

	return nil
}

func InsertDocumentsData(db *sql.DB) error {
	document := struct {
		Title        string
		FileUrl      string
		DocumentType string
	}{
		"Đề số ", "https://example.com/de", "PDF",
	}

	queryInsertDocument := `INSERT INTO documents (subjectId, title, fileUrl, documentType, views, downloads, author) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for index := 1; index < 500; index++ {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := rnd.Intn(1000) + 1
		name := gofakeit.Name()
		views := gofakeit.Number(1, 1000)
		downloads := gofakeit.Number(1, 1000)
		title := gofakeit.Book().Title

		IdStr := strconv.Itoa(randomNumber)
		_, err := db.Exec(queryInsertDocument, randomNumber%135+1, title, document.FileUrl+IdStr+".pdf", document.DocumentType, views, downloads, name)
		if err != nil {
			return fmt.Errorf("failed to insert document %s: %w", document.Title+IdStr, err)
		}
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
	InsertClassesData(db)
	InsertSubjectsData(db)
	InsertDocumentsData(db)
}
