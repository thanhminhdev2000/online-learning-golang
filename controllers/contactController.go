package controllers

import (
	"database/sql"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"

	"github.com/gin-gonic/gin"
)

// Contact godoc
// @Summary Send email contact
// @Description Send email contact
// @Tags Contact
// @Produce json
// @Param contact body models.Contact true "Send email"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error
// @Router /contact [post]
func Contact(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var content models.Contact
		if err := c.ShouldBindJSON(&content); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		err := utils.SendContactEmail(content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to send contact email"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Contact email sent successful"})
	}
}
