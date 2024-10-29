package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func createUserCommon(db *sql.DB, user *models.CreateUserRequest) (*models.UserDetail, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
	if err := db.QueryRow(query, user.Email, user.Username).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("email or username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	if user.Role == "" {
		user.Role = models.RoleUser
	}

	// Use a transaction to ensure we can get the created user's ID
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO users (email, username, fullName, password, gender, avatar, dateOfBirth, role, phoneNumber) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email, user.Username, user.FullName, hashedPassword,
		user.Gender, user.Avatar, user.DateOfBirth, user.Role, user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get created user ID: %v", err)
	}

	// Fetch the created user
	var createdUser models.UserDetail
	err = tx.QueryRow(`
		SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role, phoneNumber 
		FROM users WHERE id = ?`,
		userID).Scan(
		&createdUser.ID, &createdUser.Email, &createdUser.Username,
		&createdUser.FullName, &createdUser.Gender, &createdUser.Avatar,
		&createdUser.DateOfBirth, &createdUser.Role, &createdUser.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &createdUser, nil
}

// CreateUser godoc
// @Summary Register a new regular user
// @Description Register a new user with email, username, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 200 {object} models.CreateUserResponse
// @Failure 400 {object} models.Error
// @Failure 409 {object} models.Error
// @Router /users/ [post]
func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.CreateUserRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		// Force regular user role
		user.Role = models.RoleUser

		createdUser, err := createUserCommon(db, &user)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "already exists"):
				c.JSON(http.StatusConflict, models.Error{Error: err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, models.CreateUserResponse{
			Message: "User registered successfully",
			User:    *createdUser,
		})
	}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 200 {object} models.CreateUserResponse
// @Failure 400 {object} models.Error
// @Failure 409 {object} models.Error
// @Router /users/admin [post]
func CreateUserAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.CreateUserRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" && user.Role == "admin" {
			c.JSON(http.StatusForbidden, models.Error{Error: "Only admins are allowed to create admin users."})
			return
		}

		createdUser, err := createUserCommon(db, &user)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "already exists"):
				c.JSON(http.StatusConflict, models.Error{Error: err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, models.CreateUserResponse{
			Message: "User registered successfully",
			User:    *createdUser,
		})
	}
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information by user ID
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body models.UserDetail true "User data"
// @Success 200 {object} models.UpdateUserResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id} [put]
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		currentUserID, _ := c.Get("userId")
		currentUserRole, _ := c.Get("role")

		// Allow users to update their own info or admin to update any user
		if currentUserRole != "admin" && currentUserID != userID {
			c.JSON(http.StatusForbidden, models.Error{
				Error: "Permission denied: can only update your own profile",
			})
			return
		}

		var updateUser models.UserDetail
		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid request body: " + err.Error(),
			})
			return
		}

		// Validate email and username uniqueness
		if err := validateUniqueFields(db, userID, updateUser.Email, updateUser.Username); err != nil {
			c.JSON(http.StatusConflict, models.Error{
				Error: err.Error(),
			})
			return
		}

		// Use transaction for update
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		query := `
			UPDATE users 
			SET email = ?, 
				username = ?, 
				fullName = ?, 
				gender = ?, 
				dateOfBirth = ? 
			WHERE id = ? AND deletedAt IS NULL`

		result, err := tx.Exec(query,
			updateUser.Email,
			updateUser.Username,
			updateUser.FullName,
			updateUser.Gender,
			updateUser.DateOfBirth,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to update user",
			})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			c.JSON(http.StatusNotFound, models.Error{
				Error: "User not found or no changes made",
			})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit transaction",
			})
			return
		}

		// Fetch updated user details
		updatedUser, err := GetUserDetail(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch updated user details",
			})
			return
		}

		c.JSON(http.StatusOK, models.UpdateUserResponse{
			Message: "User updated successfully",
			User:    updatedUser,
		})
	}
}

// Password godoc
// @Summary Change user password
// @Description Change the user's password. Users can change their own password, admins can change any user's password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param password body models.PasswordUpdateRequest true "Password data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /users/{id}/password [put]
func UpdateUserPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		currentUserID, _ := c.Get("userId")
		currentUserRole, _ := c.Get("role")

		// Check permissions - allow users to change their own password or admin to change any password
		if currentUserRole != "admin" && currentUserID != userID {
			c.JSON(http.StatusForbidden, models.Error{
				Error: "Permission denied: can only update your own password",
			})
			return
		}

		var req models.PasswordUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid request body: " + err.Error(),
			})
			return
		}

		// Validate new password
		if len(req.NewPassword) < 6 {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "New password must be at least 6 characters long",
			})
			return
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		var storedPassword string
		err = tx.QueryRow("SELECT password FROM users WHERE id = ? AND deletedAt IS NULL", userID).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{
					Error: "User not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch user details",
			})
			return
		}

		// Verify current password
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{
				Error: "Current password is incorrect",
			})
			return
		}

		// Hash new password
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to process new password",
			})
			return
		}

		// Update password
		_, err = tx.Exec("UPDATE users SET password = ? WHERE id = ? AND deletedAt IS NULL",
			hashedNewPassword, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to update password",
			})
			return
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit password update",
			})
			return
		}

		c.JSON(http.StatusOK, models.Message{
			Message: "Password updated successfully",
		})
	}
}

// UpdateUserAvatar godoc
// @Summary Update user avatar
// @Description Update the avatar for a specific user. Users can update their own avatar, admins can update any user's avatar
// @Tags User
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param userId path int true "User ID"
// @Param avatar formData file true "User Avatar"
// @Success 200 {object} models.UpdateUserResponse
// @Failure 400 {object} models.Error
// @Failure 403 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id}/avatar [put]
func UpdateUserAvatar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		currentUserID, _ := c.Get("userId")
		currentUserRole, _ := c.Get("role")

		if currentUserRole != "admin" && currentUserID != userID {
			c.JSON(http.StatusForbidden, models.Error{
				Error: "Permission denied: can only update your own avatar",
			})
			return
		}

		// Validate user exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND deletedAt IS NULL)", userID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to verify user existence",
			})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, models.Error{
				Error: "User not found",
			})
			return
		}

		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "No file provided or invalid file",
			})
			return
		}

		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "File size exceeds maximum limit of 5MB",
			})
			return
		}

		if !utils.IsValidImageType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid file type. Only JPEG, PNG, and GIF are allowed",
			})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to process uploaded file",
			})
			return
		}
		defer fileContent.Close()

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		// Get old avatar URL for cleanup
		var oldAvatarURL string
		err = tx.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&oldAvatarURL)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch current avatar",
			})
			return
		}

		cld, err := utils.SetupCloudinary()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to initialize upload service",
			})
			return
		}

		// Upload new avatar
		avatarURL, err := utils.UploadImage(cld, fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to upload avatar",
			})
			return
		}

		_, err = tx.Exec("UPDATE users SET avatar = ? WHERE id = ? AND deletedAt IS NULL",
			avatarURL, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to update avatar in database",
			})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit avatar update",
			})
			return
		}

		if oldAvatarURL != "" {
			go utils.DeleteImage(cld, oldAvatarURL)
		}

		user, err := GetUserDetail(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch updated user details",
			})
			return
		}

		c.JSON(http.StatusOK, models.UpdateUserResponse{
			Message: "Avatar updated successfully",
			User:    user,
		})
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft delete a user by user ID. Only admins can delete users, and admins cannot delete their own account.
// @Tags User
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} models.Message
// @Failure 403 {object} models.Error "Permission denied or trying to delete own account"
// @Failure 404 {object} models.Error "User not found"
// @Failure 500 {object} models.Error "Server error"
// @Router /users/{id} [delete]
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, _ := c.Get("userId")
		currentUserRole, _ := c.Get("role")
		userIDToDelete := c.Param("id")

		// Check admin permission
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, models.Error{
				Error: "Permission denied: admin access required",
			})
			return
		}

		// Prevent admin from deleting their own account
		if userIDToDelete == currentUserID {
			c.JSON(http.StatusForbidden, models.Error{
				Error: "Cannot delete your own admin account",
			})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		// Check if user exists and is not already deleted
		var exists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND deletedAt IS NULL)", userIDToDelete).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to verify user existence",
			})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, models.Error{
				Error: "User not found or already deleted",
			})
			return
		}

		var avatarURL string
		err = tx.QueryRow("SELECT avatar FROM users WHERE id = ?", userIDToDelete).Scan(&avatarURL)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch user details",
			})
			return
		}

		result, err := tx.Exec(`
			UPDATE users 
			SET deletedAt = NOW(),
				email = CONCAT(email, '_deleted_', id),
				username = CONCAT(username, '_deleted_', id)
			WHERE id = ? AND deletedAt IS NULL`,
			userIDToDelete)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to delete user",
			})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to verify deletion",
			})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, models.Error{
				Error: "User not found or already deleted",
			})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit user deletion",
			})
			return
		}

		if avatarURL != "" {
			go func() {
				if cld, err := utils.SetupCloudinary(); err == nil {
					utils.DeleteImage(cld, avatarURL)
				}
			}()
		}

		c.JSON(http.StatusOK, models.Message{
			Message: "User deleted successfully",
		})
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users, with optional filters for email, username, full name, date of birth, role, and pagination.
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param data query models.UserQueryParams false "Filter"
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.Error
// @Router /users/ [get]
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != "admin" {
			c.JSON(http.StatusForbidden, models.Error{Error: "Permission denied: admin access required"})
			return
		}

		filters := map[string]string{
			"email":       c.Query("email"),
			"username":    c.Query("username"),
			"fullName":    c.Query("fullName"),
			"dateOfBirth": c.Query("dateOfBirth"),
			"role":        c.Query("role"),
		}

		// Handle pagination
		page := utils.ParseIntWithDefault(c.Query("page"), 1)
		limit := utils.ClampInt(utils.ParseIntWithDefault(c.Query("limit"), 10), 1, 100)
		offset := (page - 1) * limit

		query, countQuery, params, countParams := buildUserQuery(filters)

		var totalCount int
		if err := db.QueryRow(countQuery, countParams...).Scan(&totalCount); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to count users"})
			return
		}

		query += " LIMIT ? OFFSET ?"
		params = append(params, limit, offset)

		// Execute query
		rows, err := db.Query(query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch users"})
			return
		}
		defer rows.Close()

		users := make([]models.UserDetail, 0)
		for rows.Next() {
			var user models.UserDetail
			if err := rows.Scan(
				&user.ID,
				&user.Email,
				&user.Username,
				&user.FullName,
				&user.Gender,
				&user.Avatar,
				&user.DateOfBirth,
				&user.Role,
			); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan user data"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, models.UserResponse{
			Data: users,
			Paging: models.PagingInfo{
				Page:       page,
				Limit:      limit,
				TotalCount: totalCount,
			},
		})
	}
}

func GetUserDetail(db *sql.DB, userId string) (models.UserDetail, error) {
	row := db.QueryRow("SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE id = ? AND deletedAt IS NULL", userId)

	var user models.UserDetail
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth, &user.Role); err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} models.UserDetail
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id} [get]
func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		currentUserID, _ := c.Get("userId")
		currentUserRole, _ := c.Get("role")

		// Check permissions
		if currentUserRole == "user" && currentUserID != userID {
			c.JSON(http.StatusForbidden, models.Error{Error: "Permission denied: cannot access other user's data"})
			return
		}

		user, err := GetUserDetail(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch user details"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func buildUserQuery(filters map[string]string) (string, string, []interface{}, []interface{}) {
	baseQuery := "SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE deletedAt IS NULL"
	countQuery := "SELECT COUNT(*) FROM users WHERE deletedAt IS NULL"

	var params, countParams []interface{}

	filterFields := map[string]string{
		"email":       "email LIKE ?",
		"username":    "username LIKE ?",
		"fullName":    "fullName LIKE ?",
		"dateOfBirth": "dateOfBirth LIKE ?",
		"role":        "role = ?",
	}

	for field, value := range filters {
		if value == "" {
			continue
		}

		if sqlPattern, exists := filterFields[field]; exists {
			baseQuery += " AND " + sqlPattern
			countQuery += " AND " + sqlPattern

			if field == "role" {
				params = append(params, value)
				countParams = append(countParams, value)
			} else {
				params = append(params, "%"+value+"%")
				countParams = append(countParams, "%"+value+"%")
			}
		}
	}

	return baseQuery, countQuery, params, countParams
}

// validateUniqueFields checks if email and username are unique for other users
func validateUniqueFields(db *sql.DB, userID string, email string, username string) error {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE (email = ? OR username = ?) 
			AND id != ? 
			AND deletedAt IS NULL
		)`

	if err := db.QueryRow(query, email, username, userID).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check unique fields: %v", err)
	}

	if exists {
		return fmt.Errorf("email or username already exists")
	}

	return nil
}
