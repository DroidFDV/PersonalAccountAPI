package app

import (
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()

	user := router.Group("/user")
	user.POST("/", handler.AddUser)
	user.GET("/:id", handler.GetUserByID)
	user.PUT("/", handler.UpdateUser)

	login := router.Group("/login")
	login.POST("/", handler.Login)

	file := router.Group("/file")
	file.POST("/upload", handler.UploadFile)

	router.Use(metrics.PrometheusMiddleware())

	return router
}
