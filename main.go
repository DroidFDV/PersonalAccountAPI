package main

import (
	"github.com/gin-gonic/gin"
)

var user struct {
	id int 
	login string 
	password string 
}

var users = []user{
	{id: 1, login: "droid", password: "1"},
	{id: 2, login: "mtvy", password: "2"}
}

func main() {
	
	router := gin.Default()

	// Определяем маршрут и обработчик
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// Запускаем сервер на порту 8080
	router.Run(":8080") // По умолчанию запускается на localhost:8080
}