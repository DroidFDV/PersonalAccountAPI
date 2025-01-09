package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type user struct {
	Id       int    `json:"Id"`
	Login    string `json:"Login"`
	Password string `json:"Password"`
}

// todo: переносим в БД
// от этого отказываемся, теперь берём данные в базе
var users = []user{
	{Id: 1, Login: "droId", Password: "1"},
	{Id: 2, Login: "mtvy", Password: "2"},
}

// var conn *pgx.Conn

func main() {

	// получение conn
	// example "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres"
	// conn, err := pgx.Connect("postgre://....

	router := gin.Default()

	// todo: перенести в функции func(c *gin.Context) {...
	router.POST("/login", func(c *gin.Context) {
		var user user

		if err := c.ShouldBind(&user); err != nil {
			// todo: log.Error(errors.Wrap(err, "login ShouldBind")) | github.com/pkg/errors -> errors.Wrap(err, "login ShouldBind")
			c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId input"})
			return
		}

		// todo: проверять в базе
		// пишешь новую функцию для похода в БД
		for _, u := range users {
			if u.Login == user.Login && u.Password == user.Password {
				c.JSON(http.StatusOK, gin.H{"Id": u.Id})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	router.GET("/user/:id", func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// todo: тоже убрать в поход в базу
		for _, user := range users {
			if user.Id == id {
				c.JSON(http.StatusOK, gin.H{"user": user.Login})
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	// todo: добавить
	// POST /user
	// PUT /user редактирование передавай на изменение полностью все поля

	// Запускаем сервер на порту 8080
	router.Run(":8080") // По умолчанию запускается на localhost:8080
}
