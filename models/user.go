package models

import (
	libs "github.com/TIG/api-sqlite/helpers"
)

// LoginObject struct
type LoginObject struct {
	Username string
	Password string
}

// User struct
type User struct {
	UserID   int    `gorm:"column:UserId;primary_key;AUTO_INCREMENT;not null" json:"UserId"`
	UserName string `gorm:"column:UserName;type:nvarchar(50)"`
	Password string `gorm:"column:Password;type:nvarchar(50)"`
	RoleCode string `gorm:"column:RoleCode;type:nvarchar(50)"`
	Fullname string `gorm:"column:FullName;type:nvarchar(50)"`
	Token    string `gorm:"column:Token;type:nvarchar(50)"`
}

// TableName func
func (User) TableName() string {
	return "User"
}

// PassBodyJSONToModel func
func (u *User) PassBodyJSONToModel(JSONObject map[string]interface{}) {
	var (
		res interface{}
		val string
	)
	val, res = libs.PassValueFromJSONObjectToVariable("UserName", JSONObject)
	if res != nil {
		u.UserName = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Password", JSONObject)
	if res != nil {
		u.Password = libs.EncryptPassword(val)
	}

	val, res = libs.PassValueFromJSONObjectToVariable("RoleCode", JSONObject)
	if res != nil {
		u.RoleCode = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("FullName", JSONObject)
	if res != nil {
		u.Fullname = val
	}

	val, res = libs.PassValueFromJSONObjectToVariable("Token", JSONObject)
	if res != nil {
		u.Token = val
	}
	return
}
