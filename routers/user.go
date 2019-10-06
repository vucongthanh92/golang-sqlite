package routers

import (
	"github.com/TIG/api-sqlite/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoute func
func UserRoute(r *gin.RouterGroup) {
	r.GET("/GetAllUsers", controllers.GetAllUsers)
	r.GET("/GetUserById/:id", controllers.GetUserByID)
	r.POST("/AddUser", controllers.AddUser)
	r.PUT("/UpdateUser/:id", controllers.UpdateUser)
	r.DELETE("/DeleteUser/:id", controllers.DeleteUser)
	r.POST("/Login", controllers.Login)
	r.PUT("/Logout/:token", controllers.Logout)
	r.PUT("/ChangePassword", controllers.ChangePassword)
}
