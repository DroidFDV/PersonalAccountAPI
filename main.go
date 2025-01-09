package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type user struct {
	Id 			int 	`json:"Id"`
	Login 		string 	`json:"Login"`
	Password 	string 	`json:"Password"`
}

var users = []user{
	{Id: 1, Login: "droId", Password: "1"},
	{Id: 2, Login: "mtvy", Password: "2"},
}

func main() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		// Login := c.Query("Login")
		// Password := c.Query("Password")
		// Login, Password := c.GetRawData()
		var user user

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId input"})
			return
		}
		// для таких целей лучше использоывать карту, где ключ будет Login,
		// а значение Password, но я пока не разбирался
		// хотя у нас есть Id как ключ
		for _, u := range users {
			if u.Login == user.Login && u.Password == user.Password {
				c.JSON(http.StatusOK, gin.H{"Id": u.Id})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	router.GET("/user/:id", func(c *gin.Context) {
		// id := c.GetRawData()
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		for _, user := range users {
			if user.Id == id {
				c.JSON(http.StatusOK, gin.H{"user": user.Login})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	// Запускаем сервер на порту 8080
	router.Run(":8080") // По умолчанию запускается на localhost:8080
}