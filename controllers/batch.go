package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	sqlite "github.com/TIG/api-sqlite/database"
	libs "github.com/TIG/api-sqlite/helpers"
	"github.com/TIG/api-sqlite/models"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// GetAllBatchs API
func GetAllBatchs(c *gin.Context) {
	defer libs.RecoverError(c, "GetAllBatchs")

	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status          = 200
		responseData    = gin.H{}
		Models          []models.Batch
		arrListProducts []map[string]interface{}
	)

	db := sqlite.Connect()
	defer db.Close()
	db.Find(&Models)
	for _, model := range Models {
		productObject := structs.Map(model)
		productObject["Balance"] = model.OrderMadeQty - model.PostedQuantity
		if (model.OrderMadeQty - model.PostedQuantity) == 0 {
			productObject["Pending"] = 0
		} else {
			productObject["Pending"] = 1
		}
		arrListProducts = append(arrListProducts, productObject)
	}
	if len(arrListProducts) > 0 {
		responseData = gin.H{
			"status": status,
			"msg":    "Success",
			"data":   arrListProducts,
		}
	} else {
		responseData = gin.H{
			"status": status,
			"msg":    "Success",
			"data":   make([]string, 0),
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// GetBatchByID API
func GetBatchByID(c *gin.Context) {
	defer libs.RecoverError(c, "GetBatchByID")

	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status       = 200
		responseData = gin.H{}
		Model        models.Batch
	)
	db := sqlite.Connect()
	defer db.Close()
	ID := c.Param("id")
	resultFind := db.Where("BatchId = ?", ID).Find(&Model)
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

// AddBatch API
func AddBatch(c *gin.Context) {
	defer libs.RecoverError(c, "AddBatch")
	f, errF := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		log.Println(errF)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("-----------AddBatch----------")

	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status       = 200
		responseData = gin.H{}
		Model        models.Batch
	)
	db := sqlite.Connect()
	defer db.Close()
	body, _ := ioutil.ReadAll(c.Request.Body)
	var JSONObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &JSONObject)
	logger.Println("body: ", string(body))
	Model.PassBodyJSONToModel(JSONObject)
	resultFind := db.Where("ProductionOrderNumber = ? and ProductCode = ?", Model.ProductionOrderNumber, Model.ProductCode).First(&Model)
	if resultFind.RowsAffected > 0 {
		Model.PassBodyJSONToModel(JSONObject)
		resultSave := db.Save(&Model)
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
				"msg":    "Update success",
			}
		}
	} else {
		resultCreate := db.Create(&Model)
		if resultCreate.RowsAffected <= 0 {
			status = 500
			responseData = gin.H{
				"status": status,
				"msg":    resultCreate.Error.Error(),
			}
		} else {
			responseData = gin.H{
				"status": status,
				"data":   Model,
				"msg":    "Add success",
			}
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// UpdateBatch API
func UpdateBatch(c *gin.Context) {
	defer libs.RecoverError(c, "UpdateBatch")
	f, errF := os.OpenFile(libs.LoggerPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		log.Println(errF)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("-------------UpdateBatch-------------")

	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status       = 200
		responseData = gin.H{}
		Model        models.Batch
	)
	db := sqlite.Connect()
	defer db.Close()
	var JSONObject map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &JSONObject)
	logger.Println("body: ", string(body))
	ID := c.Param("id")
	resultFind := db.Where("BatchId = ?", ID).First(&Model)
	if resultFind.RowsAffected <= 0 {
		status = 404
		responseData = gin.H{
			"status": status,
			"msg":    "Id not found",
		}
	} else {
		Model.PassBodyJSONToModel(JSONObject)
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
				"msg":    "Update success",
			}
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// DeleteBatch API
func DeleteBatch(c *gin.Context) {
	defer libs.RecoverError(c, "DeleteBatch")

	statusToken, response := CheckTokenAPI(c)
	if statusToken != 200 {
		libs.ResponseRESTAPI(response, c, statusToken)
		return
	}

	var (
		status       = 200
		responseData = gin.H{}
		msg          string
		Model        models.Batch
	)
	db := sqlite.Connect()
	defer db.Close()
	ID := c.Param("id")
	result := db.Where("batchId = ?", ID).First(&Model)
	if result.RowsAffected <= 0 {
		status = 404
		msg = "ID not found"
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
