package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type user struct {
	id 			int 	`json:"id"`
	login 		string 	`json:"login"`
	password 	string 	`json:"password"`
}

var users = []user{
	{id: 1, login: "droid", password: "1"},
	{id: 2, login: "mtvy", password: "2"},
}

func main() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		login := c.Query("login")
		password := c.Query("password")
		// для таких целей лучше использоывать карту, где ключ будет login,
		// а значение password, но я пока не разбирался
		for _, user := range users {
			if user.login == login && user.password == password {
				c.JSON(http.StatusOK, gin.H{"id": user.id})
				return
			}
		}
		c.JSON(401, gin.H{"error": "Unauthorized"})
	})

	// Запускаем сервер на порту 8080
	router.Run(":8080") // По умолчанию запускается на localhost:8080
}