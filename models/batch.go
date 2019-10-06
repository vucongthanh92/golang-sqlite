package models

import (
	"strconv"
	"time"

	libs "github.com/TIG/api-sqlite/helpers"
)

// Batch struct
type Batch struct {
	BatchID               int        `gorm:"column:BatchId;primary_key;AUTO_INCREMENT;not null" json:"BatchId"`
	OrderDate             *time.Time `gorm:"column:OrderDate;"`
	PrinterLabel          string     `gorm:"column:PrinterLabel;type:nvarchar(50)"`
	ProductionOrderNumber string     `gorm:"column:ProductionOrderNumber;type:nvarchar(50)"`
	ProductCode           string     `gorm:"column:ProductCode;type:nvarchar(50)"`
	Batch                 string     `gorm:"column:Batch;type:nvarchar(50)"`
	OrderMadeQty          int        `gorm:"column:OrderMadeQty;"`
	PostedQuantity        int        `gorm:"column:PostedQuantity;"`
	Description           string     `gorm:"column:Description;"`
	Cost                  float64    `gorm:"column:Cost;"`
	Price                 float64    `gorm:"column:Price;"`
	RetailPrice           float64    `gorm:"column:RetailPrice;"`
}

// TableName func
func (Batch) TableName() string {
	return "Batch"
}

// PassBodyJSONToModel func
func (u *Batch) PassBodyJSONToModel(JSONObject map[string]interface{}) {
	var (
		res interface{}
		val string
	)
	val, res = libs.PassValueFromJSONObjectToVariable("OrderDate", JSONObject)
	if res != nil {
		vOrderDate, sOrderDate := time.Parse(libs.FORMATDATEFULL, val)
		if sOrderDate == nil {
			u.OrderDate = &vOrderDate
		}
	}

	val, res = libs.PassValueFromJSONObjectToVariable("PrinterLabel", JSONObject)
	if res != nil {
		u.PrinterLabel = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("ProductionOrderNumber", JSONObject)
	if res != nil {
		u.ProductionOrderNumber = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("ProductCode", JSONObject)
	if res != nil {
		u.ProductCode = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Batch", JSONObject)
	if res != nil {
		u.Batch = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("OrderMadeQty", JSONObject)
	if res != nil {
		vOrderMadeQty, sOrderMadeQty := strconv.Atoi(val)
		if sOrderMadeQty == nil {
			u.OrderMadeQty = vOrderMadeQty
		}
	}

	val, res = libs.PassValueFromJSONObjectToVariable("PostedQuantity", JSONObject)
	if res != nil {
		vPostedQuantity, sPostedQuantity := strconv.Atoi(val)
		if sPostedQuantity == nil {
			u.PostedQuantity = vPostedQuantity
		}
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Description", JSONObject)
	if res != nil {
		u.Description = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Cost", JSONObject)
	if res != nil {
		vCost, sCost := strconv.ParseFloat(val, 64)
		if sCost == nil {
			u.Cost = vCost
		}
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Price", JSONObject)
	if res != nil {
		vPrice, sPrice := strconv.ParseFloat(val, 64)
		if sPrice == nil {
			u.Price = vPrice
		}
	}

	val, res = libs.PassValueFromJSONObjectToVariable("RetailPrice", JSONObject)
	if res != nil {
		vRetailPrice, sRetailPrice := strconv.ParseFloat(val, 64)
		if sRetailPrice == nil {
			u.RetailPrice = vRetailPrice
		}
	}

	return
}
