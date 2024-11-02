package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
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

func DropCoursesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS courses;`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop courses table: %w", err)
	}
	return nil
}

func DropLessonsTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS lessons;`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop lessons table: %w", err)
	}
	return nil
}

func DropPurchasesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS purchases;`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop purchases table: %w", err)
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
        avatar VARCHAR(255) NOT NULL DEFAULT "",
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
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
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
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
		views INT DEFAULT 0,
        downloads INT DEFAULT 0,
		author VARCHAR(255) DEFAULT "admin",
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (subjectId) REFERENCES subjects(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}
	return nil
}

func CreateCoursesTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS courses (
        id INT AUTO_INCREMENT PRIMARY KEY,
        subjectId INT NOT NULL,
        title VARCHAR(255) NOT NULL,
		thumbnailUrl VARCHAR(255) NOT NULL,
        description TEXT,
        price DECIMAL(10, 2) NOT NULL,
        instructor VARCHAR(255),
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (subjectId) REFERENCES subjects(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create courses table: %w", err)
	}
	return nil
}

func CreateLessonsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS lessons (
        id INT AUTO_INCREMENT PRIMARY KEY,
        courseId INT NOT NULL,
        title VARCHAR(255) NOT NULL,
        videoUrl VARCHAR(255) NOT NULL,
        duration INT NOT NULL,
        position INT NOT NULL,
        FOREIGN KEY (courseId) REFERENCES courses(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create lessons table: %w", err)
	}
	return nil
}

func CreatePurchasesTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS purchases (
        id INT AUTO_INCREMENT PRIMARY KEY,
        userId INT NOT NULL,
        courseId INT NOT NULL,
        purchaseDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (courseId) REFERENCES courses(id) ON DELETE CASCADE
    );`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create purchases table: %w", err)
	}
	return nil
}

func InsertTestAccounts(db *sql.DB) error {
	query := `
	INSERT INTO users (email, username, fullName, password, gender, dateOfBirth, role)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	gofakeit.Seed(0)

	for index := 1; index <= 100; index++ {
		var role string
		var gender string
		if index <= 10 {
			role = "admin"
		} else {
			role = "user"
		}

		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		random := rnd.Intn(10)
		if random%2 == 0 {
			gender = "male"
		} else {
			gender = "female"
		}

		IdStr := strconv.Itoa(index)
		minDate := time.Now().AddDate(-50, 0, 0)
		maxDate := time.Now().AddDate(-12, 0, 0)
		_, err := db.Exec(query, gofakeit.Email(), role+IdStr, gofakeit.Name(), "$2a$10$3q1Qcjx7zzpb3Vs42D6YbexPA4K9pKVA9pA2T8UIo0TjccGmet10m", gender, gofakeit.DateRange(minDate, maxDate), role)

		if err != nil {
			return fmt.Errorf("failed to insert document : %w", err)
		}
	}

	return nil
}

func InsertClassesData(db *sql.DB) error {
	classes := []string{
		"Đề thi thử đại học", "Lớp 12", "Lớp 11", "Lớp 10",
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
		"Đề thi thử đại học", "Lớp 12", "Lớp 11", "Lớp 10",
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
}

func InsertDocumentsData(db *sql.DB) error {
	query := `INSERT INTO documents (subjectId, title, fileUrl, views, downloads) VALUES (?, ?, ?, ?, ?)`
	awsStorage := os.Getenv("AWS_STORAGE")

	data := []Upload{
		{1, "Đề thi thử môn Toán tỉnh Bắc Ninh"},
		{1, "Đề thi thử môn Toán tỉnh Nghệ An"},
		{1, "Đề thi thử môn Toán tỉnh Phú Yên"},
		{1, "Đề thi thử môn Toán tỉnh Quảng Trị"},
		{1, "Đề thi thử môn Toán tỉnh Thanh Hoá"},

		{2, "Đề thi thử môn Văn tỉnh Bắc Ninh"},
		{2, "Đề thi thử môn Văn tỉnh Hà Nội"},
		{2, "Đề thi thử môn Văn tỉnh Hải Dương"},
		{2, "Đề thi thử môn Văn tỉnh Kon Tum"},
		{2, "Đề thi thử môn Văn tỉnh Nghệ An"},
		{2, "Đề thi thử môn Văn tỉnh Vĩnh Phúc"},

		{3, "Đề thi thử môn Tiếng Anh tỉnh Bình Định"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Hà Nội"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Hoà Bình"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Quảng Nam"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Quảng Ninh"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Thái Nguyên"},
		{3, "Đề thi thử môn Tiếng Anh tỉnh Thanh Hoá"},

		{4, "Đề thi giữa kỳ lớp 12 môn Toán trường THPT Ngô Thời Nhiệm"},
		{4, "Đề thi giữa kỳ lớp 12 môn Toán trường THPT Nguyễn Thái Học"},
		{4, "Đề thi giữa kỳ lớp 12 môn Toán trường THPT Yên Hoà"},
		{4, "Đề thi KSCL lớp 12 môn Toán trường THPT Đoàn Kết"},
		{4, "Đề thi KSCL lớp 12 môn Toán trường THPT Đông Đậu"},
		{4, "Đề thi KSCL lớp 12 môn Toán trường THPT Nguyễn Khuyến"},

		{5, "Đề kiểm tra HKI lớp 12 môn Văn tỉnh Đồng Tháp"},
		{5, "Đề kiểm tra HKII lớp 12 môn Văn tỉnh Hà Nội"},
		{5, "Đề kiểm tra HKII lớp 12 môn Văn tỉnh Hà Tĩnh"},
		{5, "Đề kiểm tra HKII lớp 12 môn Văn tỉnh Thái Nguyên"},

		{6, "20 chuyên đề ngữ pháp Tiếng Anh ôn thi"},
		{6, "3000 từ vựng tiếng Anh thông dụng"},
		{6, "Đề thi thử THPT QG 2019 môn Tiếng Anh tỉnh Vĩnh Phúc"},
		{6, "Một số cấu trúc viết lại câu"},
		{6, "Từ vựng SGK lớp 12"},

		{7, "Đề kiểm tra Đại số chương 1"},
		{7, "Đề kiểm tra Đại số chương 2"},
		{7, "Đề kiểm tra Đại số chương 3"},
		{7, "Đề kiểm tra Hình học chương 1"},
		{7, "Đề kiểm tra Hình học chương 2"},

		{8, "Đề thi HKII lớp 11 môn Ngữ văn tỉnh Hà Nam"},
		{8, "Đề thi HKII lớp 11 môn Ngữ văn tỉnh Hà Nội"},
		{8, "Đề thi HKII lớp 11 môn Ngữ văn tỉnh Thái Bình"},

		{9, "Đề cương ôn thi HKII"},
		{9, "Tổng hợp từ vựng Unit 9+10+11"},

		{10, "Lời giải đề KT HKI - Bứt phá 9+ môn Toán"},
		{10, "Lời giải đề KT HKII - Bứt phá 9+ môn Toán"},

		{11, "7 dạng đề nghị luận thường gặp"},
		{11, "Kiến thức quan trọng Ngữ văn lớp 9"},

		{12, "Bài tập trắc nghiệm Tiếng Anh"},
		{12, "Lấy lại gốc tiếng Anh - 12 thì của động từ"},
		{12, "Tổng hợp bài tập Tiếng Anh lớp 10"},

		{13, "46 đề Toán tự luyện thi vào lớp 10"},
		{13, "Bộ đề thi giữa học kì 1 môn Toán lớp 9"},
		{13, "Chuyên đề rút gọn biểu thức"},
		{13, "Đề thi giữa học kỳ I môn Toán 9 (có đáp án)"},

		{14, "22 bài văn mẫu lớp 9"},
		{14, "Mùa xuân nho nhỏ"},
		{14, "Viếng lăng Bác"},

		{15, "Cách sử dụng câu điều kiện loại I"},
		{15, "Cách sử dụng câu điều kiện loại II"},
		{15, "Cách sử dụng câu điều kiện loại II"},
		{15, "Đề thi giữa kỳ HKI môn tiếng Anh lớp 9"},
		{15, "Mẹo làm bài tập dạng đặt câu hỏi với từ"},
		{15, "200 câu trắc nghiệm tiếng Anh lớp 9 có đáp án"},

		{16, "Các trường hợp đồng dạng của tam giác vuông"},
		{16, "Hệ thống kiến thức trọng tâm môn Toán 8"},

		{17, "Miêu tả và biểu cảm trong văn bản tự sự"},
		{17, "Phương pháp làm bài văn nghị luận"},
		{17, "Tính thống nhất về chủ đề của văn bản"},
		{17, "Tuyển tập những bài văn mẫu hay lớp 8"},
		{17, "Xây dựng đoạn văn trong văn bản"},

		{18, "Bài tập trắc nghiệm tiếng anh lớp 8"},
		{18, "Đề thi  tiếng Anh HKI"},
		{18, "Unit 10 - Communication"},

		{19, "30 đề thi học sinh giỏi toán lớp 7 có đáp án"},
		{19, "Chuyên đề về luỹ thừa số hữu tỉ"},
		{19, "Tổng hợp kiến thức quan trọng môn Toán lớp 7"},

		{20, "Ôn tập văn nghị luận Ngữ văn lớp 7"},
		{20, "Tìm hiểu chung về văn bản hành chính"},
		{20, "Tình yêu nước của nhân dân ta"},

		{21, "130 câu trắc nghiệm tiếng Anh lớp 7"},
		{21, "Bài tập ôn hè môn tiếng Anh lớp 7"},
		{21, "Đề thi khảo sát chất lượng đầu năm môn tiếng Anh lớp 7"},

		{22, "Bộ đề kiểm tra 15 phút chương 1 Số học lớp 6"},
		{22, "Bộ đề kiểm tra 15 phút chương 2 Số học lớp 6"},
		{22, "Tổng hợp kiến thức môn Toán lớp 6"},
		{22, "Đề thi khảo sát chất lượng đầu năm môn Toán lớp 6"},
		{22, "40 đề kiểm tra khảo sát chất lượng đầu năm môn Toán 6"},

		{23, "Đoạn văn tự giới thiệu về bản thân"},
		{23, "Tả cảnh buổi sáng mùa xuân trên quê hương"},

		{24, "185 câu chia động từ tiếng Anh lớp 6-7"},
		{24, "Bài tập thì hiện tại đơn & hiện tại tiếp diễn"},
		{24, "Văn mẫu về tác hại của việc chơi game"},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, row := range data {
		views := rnd.Intn(5000) + 1500
		downloads := rnd.Intn(1000) + 200

		_, err := db.Exec(query, row.SubjectId, row.Title, awsStorage+"pdfs/"+row.Title+".pdf", views, downloads)
		if err != nil {
			return fmt.Errorf("failed to insert class %s: %w", row.Title, err)
		}
	}

	return nil
}

func InsertCoursesData(db *sql.DB) error {
	cloudinaryStorage := os.Getenv("CLOUDINARY_STORAGE")

	query := `INSERT INTO courses (subjectId, title, thumbnailUrl, description, price, instructor) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query,
		4,
		"Giải đề thi THPT Quốc gia bằng máy tính Casio",
		cloudinaryStorage+"image/upload/v1730551590/images/jxicfhtstqu1c0xcvsrr.jpg",
		"Giải đề thi THPT Quốc gia bằng máy tính Casio là phương pháp giúp học sinh giải nhanh các bài toán trắc nghiệm Toán. Thông qua việc sử dụng các chức năng của máy tính Casio, học sinh có thể giải quyết các dạng bài từ cơ bản đến nâng cao một cách hiệu quả và tiết kiệm thời gian. Hướng dẫn này sẽ cung cấp các mẹo và ví dụ minh họa chi tiết để hỗ trợ việc ôn luyện.",
		2000000,
		"Nguyễn Thành Minh")

	if err != nil {
		return fmt.Errorf("failed to insert course: %w", err)
	}

	return nil
}

type CreateLesson struct {
	CourseID int
	Title    string
	VideoURL string
	Duration int
	Position int
}

func InsertLessonsData(db *sql.DB) error {
	cloudinaryStorage := os.Getenv("CLOUDINARY_STORAGE")

	query := `INSERT INTO lessons (courseId, title, videoUrl, duration, position) VALUES (?, ?, ?, ?, ?)`

	data := []CreateLesson{
		{1, "[TOÁN 12] - Casio Tích phân tham số a, b (câu 1)", "video/upload/v1730552349/videos/ufpeid04iixmfukhqs0h.mp4", 225, 1},
		{1, "[THPT QUỐC GIA] - Giải đề thi bằng máy tính casio TOÁN 12 (PHẦN 1)", "video/upload/v1730554148/videos/vgsw5axuokninr2ehqee.mp4", 313, 2},
		{1, "[THPT Quốc gia] - Tích phân chứa tham số (Câu 1)", "video/upload/v1730554326/videos/bbu09gbprwgxfojxjvb1.mp4", 211, 3},
		{1, "[THPT Quốc gia] - Tích phân chứa tham số (Câu 2)", "video/upload/v1730554424/videos/wbj1sw6qbtsxjvmghkkg.mp4", 240, 4},
		{1, "[Ôn thi THPT] - Giải đề thi số 102 bằng máy tính casio - PHẦN 1", "video/upload/v1730554424/videos/wbj1sw6qbtsxjvmghkkg.mp4", 677, 5},
		{1, "[Ôn thi THPT] - Giải đề thi số 102 bằng máy tính casio - PHẦN 2", "video/upload/v1730554653/videos/jiyhhyifhau3anm2buzj.mp4", 2882, 6},
		{1, "[Thi THPT] Giải đề Minh họa bằng máy tính casio (Đề 110 - Phần 1)", "video/upload/v1730554982/videos/lf10zrhnqpckqs6cy6vk.mp4", 1384, 7},
		{1, "[Thi THPT] Giải đề Minh họa bằng máy tính casio (Đề 110 - Phần 2)", "video/upload/v1730555002/videos/t7d14a7yrhajgjkqd8mx.mp4", 1937, 8},
	}

	for _, row := range data {
		_, err := db.Exec(query, row.CourseID, row.Title, cloudinaryStorage+row.VideoURL, row.Duration, row.Position)
		if err != nil {
			return fmt.Errorf("failed to insert class %s: %w", row.Title, err)
		}
	}

	return nil
}

func ResetDataBase(db *sql.DB) error {

	if err := DropPurchasesTable(db); err != nil {
		return err
	}
	if err := DropLessonsTable(db); err != nil {
		return err
	}
	if err := DropCoursesTable(db); err != nil {
		return err
	}

	if err := DropDocumentsTable(db); err != nil {
		return err
	}
	if err := DropSubjectsTable(db); err != nil {
		return err
	}
	if err := DropClassesTable(db); err != nil {
		return err
	}

	if err := DropResetPasswordTokensTable(db); err != nil {
		return err
	}
	if err := DropUsersTable(db); err != nil {
		return err
	}

	// ----------------

	if err := CreateUsersTable(db); err != nil {
		return err
	}
	if err := CreateResetPasswordTokensTable(db); err != nil {
		return err
	}

	if err := CreateClassesTable(db); err != nil {
		return err
	}
	if err := CreateSubjectsTable(db); err != nil {
		return err
	}
	if err := CreateDocumentsTable(db); err != nil {
		return err
	}

	if err := CreateCoursesTable(db); err != nil {
		return err
	}
	if err := CreateLessonsTable(db); err != nil {
		return err
	}
	if err := CreatePurchasesTable(db); err != nil {
		return err
	}

	// ----------------

	if err := InsertTestAccounts(db); err != nil {
		return err
	}
	if err := InsertClassesData(db); err != nil {
		return err
	}
	if err := InsertSubjectsData(db); err != nil {
		return err
	}
	if err := InsertDocumentsData(db); err != nil {
		return err
	}
	if err := InsertCoursesData(db); err != nil {
		return err
	}
	if err := InsertLessonsData(db); err != nil {
		return err
	}

	return nil
}
