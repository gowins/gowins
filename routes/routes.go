package routes

import (
	"gowins/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Define routes
	r.POST("/ping", controllers.Ping)
}
