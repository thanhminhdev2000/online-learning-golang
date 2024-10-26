package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	awsSetup "online-learning-golang/aws"
	"online-learning-golang/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClassesWithSubjects(db *sql.DB) ([]models.ClassWithSubjects, error) {
	classQuery := `
        SELECT c.id, c.name, COUNT(d.id) as documentCount
        FROM classes c
        LEFT JOIN subjects s ON c.id = s.classId
        LEFT JOIN documents d ON s.id = d.subjectId
        GROUP BY c.id;
    `
	rows, err := db.Query(classQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	defer rows.Close()

	var classList []models.ClassWithSubjects

	for rows.Next() {
		var class models.ClassWithSubjects
		if err := rows.Scan(&class.ClassId, &class.ClassName, &class.Count); err != nil {
			return nil, err
		}

		subjectQuery := `
            SELECT s.id, s.name, COUNT(d.id) as documentCount
            FROM subjects s
            LEFT JOIN documents d ON s.id = d.subjectId
            WHERE s.classId = ?
            GROUP BY s.id;
        `
		subjectRows, err := db.Query(subjectQuery, class.ClassId)
		if err != nil {
			return nil, fmt.Errorf("failed to get subjects for classId %d: %w", class.ClassId, err)
		}
		defer subjectRows.Close()

		var subjects []models.SubjectId
		for subjectRows.Next() {
			var subject models.SubjectId
			if err := subjectRows.Scan(&subject.SubjectId, &subject.SubjectName, &subject.Count); err != nil {
				return nil, err
			}
			subjects = append(subjects, subject)
		}

		class.Subjects = subjects
		classList = append(classList, class)
	}

	return classList, nil
}

// GetListClassesWithSubjects godoc
// @Summary List of classes with their subjects and document counts
// @Description List of classes with their subjects and document counts
// @Tags Document
// @Success 200 {array} models.ClassWithSubjects
// @Failure 500 {object} models.Error
// @Router /documents/subjects [get]
func GetListClassesWithSubjects(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		classesWithSubjects, err := GetClassesWithSubjects(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, classesWithSubjects)
	}
}

// UploadDocument godoc
// @Summary Upload document file
// @Description Upload document file
// @Tags Document
// @Security BearerAuth
// @Param subjectId formData int true "Subject ID"
// @Param title formData string true "Document title"
// @Param author formData string false "Document author"
// @Param file formData file true "File to upload"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error
// @Router /documents/upload [post]
func UploadDocument(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var document models.UploadRequest
		if err := c.ShouldBind(&document); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
			return
		}
		defer file.Close()

		url, err := awsSetup.UploadPDF(file, fileHeader)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		query := "INSERT INTO documents (subjectId, title, fileUrl) VALUES (?, ?, ?)"
		_, err = db.Exec(query, document.SubjectId, document.Title, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload document"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Upload document successful"})
	}
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete a document by document ID
// @Tags Document
// @Security BearerAuth
// @Param documentId path int true "Document ID"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error
// @Router /documents/{documentId} [delete]
func DeleteDocument(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this document"})
			return
		}

		documentId := c.Param("documentId")
		var fileUrl string
		query := "SELECT fileUrl FROM documents WHERE id = ?"

		row := db.QueryRow(query, documentId)
		if err := row.Scan(&fileUrl); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
			return
		}

		query = "DELETE FROM documents WHERE id = ?"
		_, err := db.Exec(query, documentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
			return
		}

		err = awsSetup.DeletePDF(fileUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Document deleted successfully"})
	}
}

// GetDocuments godoc
// @Summary Lấy danh sách tài liệu
// @Description Trả về danh sách các tài liệu, có thể lọc theo `subjectId` và `title`. Giới hạn số lượng tài liệu trả về bằng tham số `limit`.
// @Tags Document
// @Param limit query int false "Giới hạn số lượng tài liệu trả về" default(40) max(100)
// @Param subjectId query int false "ID của môn học"
// @Param title query string false "Tiêu đề của tài liệu (tìm kiếm bằng LIKE)"
// @Success 200 {array} models.DocumentsResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /documents/ [get]
func GetDocuments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.Query("limit")
		limit := 40
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if limit > 100 {
				limit = 100
			}
			if err != nil || limit <= 0 {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid limit parameter. Must be a positive integer."})
				return
			}
		}

		// Nhận các query parameters cho subjectId và title
		subjectIdStr := c.Query("subjectId")
		title := c.Query("title")

		// Xây dựng câu truy vấn cơ bản
		query := `
			SELECT 
				d.id, 
				CONCAT(s.name, ' - ', c.name) AS category, 
				d.title, 
				d.fileUrl, 
				d.documentType, 
				d.views, 
				d.downloads, 
				d.author 
			FROM documents d
			JOIN subjects s ON d.subjectId = s.id
			JOIN classes c ON s.classId = c.id
		`

		// Xây dựng các điều kiện WHERE
		var conditions []string
		var args []interface{}

		if subjectIdStr != "" {
			conditions = append(conditions, "d.subjectId = ?")
			args = append(args, subjectIdStr)
		}
		if title != "" {
			conditions = append(conditions, "d.title LIKE ?")
			args = append(args, "%"+title+"%")
		}

		// Nếu có điều kiện thì thêm vào câu truy vấn
		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		// Thêm phần ORDER BY và LIMIT
		query += " ORDER BY d.views DESC LIMIT ?"
		args = append(args, limit)

		// Thực hiện câu truy vấn với các tham số
		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve documents"})
			return
		}
		defer rows.Close()

		var documents []models.DocumentsResponse
		for rows.Next() {
			var doc models.DocumentsResponse
			if err := rows.Scan(&doc.ID, &doc.Category, &doc.Title, &doc.FileUrl, &doc.DocumentType, &doc.Views, &doc.Downloads, &doc.Author); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to parse document data"})
				return
			}
			documents = append(documents, doc)
		}

		if len(documents) == 0 {
			documents = []models.DocumentsResponse{}
		}

		c.JSON(http.StatusOK, documents)
	}
}