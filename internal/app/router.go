package app

import (
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()

	router.POST("/user", handler.AddUser)
	router.POST("/login", handler.Login)
	router.GET("/user/:id", handler.GetUserByID)
	router.PUT("/user", handler.UpdateUser)
	router.POST("/file/upload", handler.UploadFile)

	router.Use(metrics.PrometheusMiddleware())

	return router
}
