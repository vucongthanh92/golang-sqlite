package libs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LOGGERPATH var
const LOGGERPATH = "logs/logger.log"

// FORMATDATEFULL var
const FORMATDATEFULL = "2006-01-02T15:04:05Z"

// LoggerPath func
func LoggerPath() string {
	array := strings.Split(LOGGERPATH, ".")
	return array[0] + time.Now().Format("-01-02-2006") + "." + array[1]
}

// EncryptPassword func
func EncryptPassword(pass string) string {
	hash := md5.Sum([]byte(pass))
	return hex.EncodeToString(hash[:])
}

// RecoverError func
func RecoverError(c *gin.Context, FuncName string) {
	if r := recover(); r != nil {
		f, err := os.OpenFile(LoggerPath(),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		logger := log.New(f, "", log.LstdFlags)
		logger.Println("----------" + FuncName + "-----------")
		logger.Println("Error: ", r)
		responseData := gin.H{
			"status": 500,
			"msg":    r,
		}
		c.JSON(500, responseData)
		return
	}
}

// ResponseRESTAPI func
func ResponseRESTAPI(responseData gin.H, c *gin.Context, status int) {
	ResponseType := c.Request.Header.Get("ResponseType")
	if ResponseType == "application/xml" {
		c.XML(status, responseData)
	} else {
		c.JSON(status, responseData)
	}
}

// PassValueFromJSONObjectToVariable func
func PassValueFromJSONObjectToVariable(FieldName string, JSONObject map[string]interface{}) (string, interface{}) {
	valPostFromObject := JSONObject[FieldName]
	valPostForm := fmt.Sprintf("%v", valPostFromObject)
	if valPostFromObject == nil || valPostForm == "" {
		valPostFormObjectLower := JSONObject[strings.ToLower(FieldName)]
		valPostFormLower := fmt.Sprintf("%v", valPostFormObjectLower)
		return valPostFormLower, valPostFormObjectLower
	}
	return valPostForm, valPostFromObject
}

// GenerateTokenAPI func
func GenerateTokenAPI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
