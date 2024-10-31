package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	awsUtils "online-learning-golang/aws"
	"online-learning-golang/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClasses(db *sql.DB) ([]models.Class, error) {
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

	var classList []models.Class

	for rows.Next() {
		var class models.Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Count); err != nil {
			return nil, err
		}

		subjectQuery := `
            SELECT s.id, s.name, COUNT(d.id) as documentCount
            FROM subjects s
            LEFT JOIN documents d ON s.id = d.subjectId
            WHERE s.classId = ?
            GROUP BY s.id;
        `
		subjectRows, err := db.Query(subjectQuery, class.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subjects for classId %d: %w", class.ID, err)
		}
		defer subjectRows.Close()

		var subjects []models.Subject
		for subjectRows.Next() {
			var subject models.Subject
			if err := subjectRows.Scan(&subject.ID, &subject.Name, &subject.Count); err != nil {
				return nil, err
			}
			subjects = append(subjects, subject)
		}

		class.Subjects = subjects
		classList = append(classList, class)
	}

	return classList, nil
}

// GetListClass godoc
// @Summary List of classes with their subjects and document counts
// @Description List of classes with their subjects and document counts
// @Tags Document
// @Success 200 {array} models.Class
// @Failure 500 {object} models.Error
// @Router /documents/subjects [get]
func GetListClass(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		class, err := GetClasses(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, class)
	}
}

// CreateDocument godoc
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
// @Router /documents/ [post]
func CreateDocument(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var document models.CreateDocument
		if err := c.ShouldBind(&document); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid file upload"})
			return
		}
		defer file.Close()

		url, err := awsUtils.UploadPDF(file, fileHeader)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload file"})
			return
		}

		query := "INSERT INTO documents (subjectId, title, fileUrl) VALUES (?, ?, ?)"
		_, err = db.Exec(query, document.SubjectID, document.Title, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create document"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Create document successful"})
	}
}

// UpdateDocument godoc
// @Summary Update a document, including replacing its file
// @Description Update a document's information and optionally replace its file by document ID
// @Tags Document
// @Security BearerAuth
// @Param id path int true "Document ID"
// @Param title formData string false "Document title"
// @Param author formData string false "Document author"
// @Param views formData int false "Number of views"
// @Param downloads formData int false "Number of downloads"
// @Param file formData file false "File to replace the existing document file"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 403 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /documents/{id} [put]
func UpdateDocument(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, models.Error{Error: "You do not have permission to update this document"})
			return
		}

		documentID := c.Param("id")
		var updateFields []string
		var args []interface{}

		if title := c.PostForm("title"); title != "" {
			updateFields = append(updateFields, "title = ?")
			args = append(args, title)
		}

		if author := c.PostForm("author"); author != "" {
			updateFields = append(updateFields, "author = ?")
			args = append(args, author)
		}

		if viewsStr := c.PostForm("views"); viewsStr != "" {
			views, err := strconv.Atoi(viewsStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid views parameter"})
				return
			}
			updateFields = append(updateFields, "views = ?")
			args = append(args, views)
		}

		if downloadsStr := c.PostForm("downloads"); downloadsStr != "" {
			downloads, err := strconv.Atoi(downloadsStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid downloads parameter"})
				return
			}
			updateFields = append(updateFields, "downloads = ?")
			args = append(args, downloads)
		}

		file, fileHeader, err := c.Request.FormFile("file")
		if err == nil {
			var oldFileUrl string

			query := "SELECT fileUrl FROM documents WHERE id = ?"
			err = db.QueryRow(query, documentID).Scan(&oldFileUrl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve existing document file"})
				return
			}

			err = awsUtils.DeletePDF(oldFileUrl)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete old document file"})
				return
			}

			newFileUrl, err := awsUtils.UploadPDF(file, fileHeader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload new document file"})
				return
			}

			updateFields = append(updateFields, "fileUrl = ?")
			args = append(args, newFileUrl)
		}

		if len(updateFields) == 0 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No fields to update"})
			return
		}

		query := fmt.Sprintf("UPDATE documents SET %s WHERE id = ?", strings.Join(updateFields, ", "))
		args = append(args, documentID)
		_, err = db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update document"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Document updated successfully"})
	}
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete a document by document ID
// @Tags Document
// @Security BearerAuth
// @Param id path int true "Document ID"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error
// @Router /documents/{id} [delete]
func DeleteDocument(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, models.Error{Error: "You do not have permission to delete this document"})
			return
		}

		documentID := c.Param("id")
		var fileUrl string
		query := "SELECT fileUrl FROM documents WHERE id = ?"

		row := db.QueryRow(query, documentID)
		if err := row.Scan(&fileUrl); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
			return
		}

		query = "DELETE FROM documents WHERE id = ?"
		_, err := db.Exec(query, documentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
			return
		}

		// Prevent delete test data
		// err = awsUtils.DeletePDF(fileUrl)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete document"})
		// 	return
		// }

		c.JSON(http.StatusOK, models.Message{Message: "Document deleted successfully"})
	}
}

// GetDocuments godoc
// @Summary Retrieve document list
// @Description Returns a list of documents, which can be filtered by `subjectId` and `title`. Limits the number of returned documents using the `limit` parameter.
// @Tags Document
// @Param limit query int false "Limit the number of documents returned" default(40) max(100)
// @Param subjectId query int false "Subject ID"
// @Param title query string false "Document title (searched using LIKE)"
// @Success 200 {array} models.Document
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

		subjectIdStr := c.Query("subjectId")
		title := c.Query("title")

		query := `
			SELECT 
				d.id,
				c.id,
				s.id, 
				CONCAT(s.name, ' - ', c.name) AS category, 
				d.title, 
				d.fileUrl, 
				d.views, 
				d.downloads, 
				d.author 
			FROM documents d
			JOIN subjects s ON d.subjectId = s.id
			JOIN classes c ON s.classId = c.id
		`

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

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		query += " ORDER BY d.views DESC LIMIT ?"
		args = append(args, limit)

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve documents"})
			return
		}
		defer rows.Close()

		var documents []models.Document
		for rows.Next() {
			var doc models.Document
			if err := rows.Scan(&doc.ID, &doc.ClassID, &doc.SubjectID, &doc.Category, &doc.Title, &doc.FileURL, &doc.Views, &doc.Downloads, &doc.Author); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to parse document data"})
				return
			}
			documents = append(documents, doc)
		}

		if len(documents) == 0 {
			documents = []models.Document{}
		}

		c.JSON(http.StatusOK, documents)
	}
}
