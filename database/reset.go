package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

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
		author VARCHAR(255) DEFAULT "admin",
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
	INSERT INTO users (email, username, fullName, password, gender, dateOfBirth, role)
	VALUES (?, ?, ?, '$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m', 'male', '1999-09-09', ?)
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
		"Lớp 9", "Lớp 8", "Lớp 7", "Lớp 6",
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
		"Toán", "Ngữ Văn", "Tiếng Anh",
	}

	queryGetClassId := `SELECT id FROM classes WHERE name = ?`
	queryInsertSubject := `INSERT INTO subjects (classId, name) VALUES (?, ?)`

	classes := []string{
		"Đề thi thử đại học", "Lớp 12", "Lớp 11", "Lớp 10", "Thi vào lớp 10",
		"Lớp 9", "Lớp 8", "Lớp 7", "Lớp 6",
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

type Upload struct {
	SubjectId int
	Title     string
	FileUrl   string
}

func InsertDocumentsData(db *sql.DB) error {
	query := `INSERT INTO documents (subjectId, title, fileUrl, views, downloads) VALUES (?, ?, ?, ?, ?)`

	data := []Upload{
		{1, "Đề thi thử môn Toán tỉnh Bắc Ninh", "https://online-learning-aws.s3.us-east-1.amazonaws.com/pdfs/Đề thi thử môn Toán tỉnh Bắc Ninh.pdf"},
		{1, "Đề thi thử môn Toán tỉnh Nghệ An", "https://online-learning-aws.s3.us-east-1.amazonaws.com/pdfs/Đề thi thử môn Toán tỉnh Nghệ An.pdf"},
		{1, "Đề thi thử môn Toán tỉnh Phú Yên", "https://online-learning-aws.s3.us-east-1.amazonaws.com/pdfs/Đề thi thử môn Toán tỉnh Phú Yên.pdf"},
		{1, "Đề thi thử môn Toán tỉnh Quảng Trị", "https://online-learning-aws.s3.us-east-1.amazonaws.com/pdfs/Đề thi thử môn Toán tỉnh Quảng Trị.pdf"},
		{1, "Đề thi thử môn Toán tỉnh Thanh Hoá", "https://online-learning-aws.s3.us-east-1.amazonaws.com/pdfs/Đề thi thử môn Toán tỉnh Thanh Hoá.pdf"},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, row := range data {
		views := rnd.Intn(5000) + 1500
		downloads := rnd.Intn(1000) + 200

		_, err := db.Exec(query, row.SubjectId, row.Title, row.FileUrl, views, downloads)
		if err != nil {
			return fmt.Errorf("failed to insert class %s: %w", row.Title, err)
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
