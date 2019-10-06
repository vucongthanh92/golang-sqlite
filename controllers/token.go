package controllers

import (
	"strings"

	sqlite "github.com/TIG/api-sqlite/database"
	"github.com/TIG/api-sqlite/models"
	"github.com/gin-gonic/gin"
)

// CheckTokenAPI func
func CheckTokenAPI(c *gin.Context) (int, gin.H) {
	var (
		status = 200
		user   models.User
		errors = make([]string, 0)
	)
	db := sqlite.Connect()
	defer db.Close()
	userid := c.Request.Header.Get("userid")
	token := c.Request.Header.Get("token")
	if userid == "" {
		status = 422
		errors = append(errors, "UserId is required")
	}
	if token == "" {
		status = 422
		errors = append(errors, "Token is required")
	}
	if status == 200 {
		result := db.Where("UserId = ?", userid).First(&user)
		if result.RowsAffected <= 0 {
			status = 404
			errors = append(errors, "userid not found")
		} else {
			if strings.TrimSpace(user.Token) != token {
				status = 401
				errors = append(errors, "token invalid")
			}
		}
	}
	return status, gin.H{"status": status, "msg": errors}
}
