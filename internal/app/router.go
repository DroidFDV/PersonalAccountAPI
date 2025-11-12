package app

import (
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()
	//NOTE: важен порядок
	router.Use(metrics.PrometheusMiddleware())

	user := router.Group("/user")
	user.POST("/", handler.AddUser)
	user.GET("/:id", handler.GetUserByID)
	user.PUT("/", handler.UpdateUser)

	login := router.Group("/login")
	login.POST("/", handler.Login)

	file := router.Group("/file")
	file.POST("/upload/:id", handler.UploadFile)

	return router
}
