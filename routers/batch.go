package routers

import (
	"github.com/TIG/api-sqlite/controllers"
	"github.com/gin-gonic/gin"
)

// BatchRoute func
func BatchRoute(r *gin.RouterGroup) {
	r.GET("/GetAllBatchs", controllers.GetAllBatchs)
	r.GET("/GetBatchById/:id", controllers.GetBatchByID)
	r.POST("/AddBatch", controllers.AddBatch)
	r.PUT("/UpdateBatch/:id", controllers.UpdateBatch)
	r.DELETE("/DeleteBatch/:id", controllers.DeleteBatch)
}
