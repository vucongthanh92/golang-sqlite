package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	sqlite "github.com/TIG/api-sqlite/database"
	libs "github.com/TIG/api-sqlite/helpers"
	"github.com/TIG/api-sqlite/models"
	"github.com/gin-gonic/gin"
)

// Login API
func Login(c *gin.Context) {
	defer libs.RecoverError(c, "Login")
	f, err := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("-----------------Login API------------------")
	var (
		status       = 200
		responseData = gin.H{}
		user         models.User
		loginObject  models.LoginObject
		errors       = make([]string, 0)
	)

	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &loginObject)
	logger.Println("body: ", string(body))
	logger.Println("loginObject: ", loginObject)
	username := loginObject.Username
	password := loginObject.Password
	db := sqlite.Connect()
	defer db.Close()
	if username == "" {
		status = 422
		errors = append(errors, "User is required")
	}
	if password == "" {
		status = 422
		errors = append(errors, "Password is required")
	}
	if status != 200 {
		responseData = gin.H{
			"status": status,
			"msg":    errors,
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}
	passwordEncrypt := libs.EncryptPassword(password)
	result := db.Where("Username = ? and Password = ?", username, passwordEncrypt).First(&user)
	if result.RowsAffected <= 0 {
		status = 401
		responseData = gin.H{
			"status": status,
			"msg":    "Invalid password or username",
		}
	} else {
		user.Token = libs.GenerateTokenAPI()
		resultSave := db.Save(&user)
		if resultSave.Error == nil {
			responseData = gin.H{
				"status": status,
				"msg":    "Success",
				"data":   user,
			}
		} else {
			status = 500
			responseData = gin.H{
				"status": status,
				"msg":    resultSave.Error.Error(),
			}
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// Logout API
func Logout(c *gin.Context) {
	defer libs.RecoverError(c, "Logout")
	f, err := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("------------------Logout API-------------------")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}
	var (
		status       = 200
		responseData = gin.H{}
		models       []models.User
	)
	db := sqlite.Connect()
	defer db.Close()
	Token := c.Param("token")
	db.Where("Token = ?", Token).Find(&models)
	for _, v := range models {
		v.Token = ""
		db.Save(&v)
	}
	responseData = gin.H{
		"status": status,
		"msg":    "Success",
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// GetAllUsers API
func GetAllUsers(c *gin.Context) {
	defer libs.RecoverError(c, "GetAllUsers")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}
	var (
		status       = 200
		responseData = gin.H{}
		Models       []models.User
	)
	db := sqlite.Connect()
	defer db.Close()
	db.Find(&Models)
	responseData = gin.H{
		"status": status,
		"msg":    "Success",
		"data":   Models,
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// GetUserByID API
func GetUserByID(c *gin.Context) {
	defer libs.RecoverError(c, "GetUserByID")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}
	var (
		status       = 200
		responseData = gin.H{}
		Model        models.User
	)
	db := sqlite.Connect()
	defer db.Close()
	UserID := c.Param("id")
	resultFind := db.Where("UserID = ?", UserID).Find(&Model)
	if resultFind.RowsAffected > 0 {
		responseData = gin.H{
			"status": status,
			"msg":    "Success",
			"data":   Model,
		}
	} else {
		status = 404
		responseData = gin.H{
			"status": status,
			"msg":    "Id not found",
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// AddUser API
func AddUser(c *gin.Context) {
	defer libs.RecoverError(c, "AddUser")
	f, errF := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		log.Println(errF)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("---------------AddUser----------------")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}
	var (
		status       = 200
		responseData = gin.H{}
		Model        models.User
	)

	db := sqlite.Connect()
	defer db.Close()
	body, _ := ioutil.ReadAll(c.Request.Body)
	var JSONObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &JSONObject)
	logger.Println("body: ", string(body))
	Model.PassBodyJSONToModel(JSONObject)
	resultFind := db.Where("UserName = ?", Model.UserName).First(&models.User{})
	if resultFind.RowsAffected <= 0 {
		resultSave := db.Create(&Model)
		if resultSave.RowsAffected <= 0 {
			status = 500
			responseData = gin.H{
				"status": status,
				"msg":    resultSave.Error.Error(),
			}
		} else {
			responseData = gin.H{
				"status": status,
				"data":   Model,
				"msg":    "Add success",
			}
		}
	} else {
		status = 422
		responseData = gin.H{
			"status": status,
			"msg":    "UserName exist",
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// UpdateUser API
func UpdateUser(c *gin.Context) {
	defer libs.RecoverError(c, "UpdateUser")
	f, errF := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		log.Println(errF)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("-----------UpdateUser-------------")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}
	var (
		status       = 200
		responseData = gin.H{}
		Model        models.User
	)
	db := sqlite.Connect()
	defer db.Close()
	var JSONObject map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &JSONObject)
	logger.Println("body: ", string(body))
	UserID := c.Param("id")
	resultFind := db.Where("UserID = ?", UserID).First(&Model)
	if resultFind.RowsAffected <= 0 {
		status = 404
		responseData = gin.H{
			"status": status,
			"msg":    "Id not found",
		}
	} else {
		Model.PassBodyJSONToModel(JSONObject)
		resultFind := db.Where("UserName = ? AND UserId <> ?", Model.UserName, UserID).First(&models.User{})
		if resultFind.RowsAffected <= 0 {
			resultSave := db.Save(&Model)
			if resultSave.RowsAffected <= 0 && resultSave.Error != nil {
				status = 500
				responseData = gin.H{
					"status": status,
					"msg":    resultSave.Error.Error(),
				}
			} else {
				responseData = gin.H{
					"status": status,
					"data":   Model,
					"msg":    "Update Success",
				}
			}
		} else {
			status = 422
			responseData = gin.H{
				"status": status,
				"msg":    "UserName exist",
			}
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// DeleteUser API
func DeleteUser(c *gin.Context) {
	defer libs.RecoverError(c, "DeleteUser")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status       = 200
		responseData = gin.H{}
		msg          string
		Model        models.User
	)
	db := sqlite.Connect()
	defer db.Close()
	UserID := c.Param("id")
	result := db.Where("UserId = ?", UserID).First(&Model)
	if result.RowsAffected <= 0 {
		status = 404
		msg = "Id not found"
	} else {
		result := db.Delete(&Model)
		if result.RowsAffected <= 0 {
			status = 500
			msg = result.Error.Error()
		} else {
			msg = "Delete success"
		}
	}
	responseData = gin.H{
		"status": status,
		"msg":    msg,
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// ChangePassword API
func ChangePassword(c *gin.Context) {
	defer libs.RecoverError(c, "ChangePassword")
	f, errF := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		log.Println(errF)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("---------------ChangePassword------------------")
	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status            = 200
		responseData      = gin.H{}
		msg               string
		arrayPasswordPost []models.LoginObject
		errors            = make([]map[string]string, 0)
		data              []interface{}
	)
	db := sqlite.Connect()
	defer db.Close()
	tx := db.Begin()
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &arrayPasswordPost)
	logger.Println("Body: ", string(body))
	logger.Println("arrayPasswordPost: ", arrayPasswordPost)
	for i, v := range arrayPasswordPost {
		if v.Username == "" {
			status = 422
			errors = append(errors, map[string]string{"username[" + strconv.Itoa(i) + "]": "Username is required"})
		}
		if v.Password == "" {
			status = 422
			errors = append(errors, map[string]string{"username[" + strconv.Itoa(i) + "]": "Password is required"})
		}
	}
	if len(arrayPasswordPost) > 0 {
		for i, v := range arrayPasswordPost {
			if status != 200 {
				break
			}
			var userModel models.User
			resultFindUser := tx.Where("UserName = ?", v.Username).First(&userModel)
			if resultFindUser.RowsAffected <= 0 {
				status = 404
				errors = append(errors, map[string]string{"username[" + strconv.Itoa(i) + "]": "User is not found"})
			} else {
				userModel.Password = libs.EncryptPassword(v.Password)
				resultUpdateUser := tx.Save(&userModel)
				if resultUpdateUser.RowsAffected <= 0 {
					status = 500
					errors = append(errors, map[string]string{"username[" + strconv.Itoa(i) + "]": resultUpdateUser.Error.Error()})
				} else {
					data = append(data, userModel)
				}
			}
		}
	}

	if status == 200 {
		tx.Commit()
		msg = "Success"
		responseData = gin.H{
			"status": status,
			"msg":    msg,
			"data":   data,
		}
	} else {
		tx.Rollback()
		msg = "Error"
		responseData = gin.H{
			"status": status,
			"msg":    msg,
			"errors": errors,
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}
