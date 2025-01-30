package handler

import (
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/usecase"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// Не сходиться с идеей, что мы хотим не изменять handler
// При каждом новом provider нам надо обновлять handler
// новое поле newProvider, новые методы, где изменяется
// только объект, метод которого вызываем
// Также по сути интерфейс не используется

// type FooInterface interface {
// 	DoSmth()
// }

// type Goo struct {
// 	fields fieldA
// }

// func (*Goo) DoSmth() {}

// type Hoo struct {
// 	fields fieldB
// }

// func (*Hoo) DoSmth() {}

// func foo(inter *FooInterface) {
// 	DoSmth()
// }

type Handle struct {
	userProvider  usecase.UserUsecase
	cacheProvider cache.CacheDecorator
}

func NewHandle(conn *pgx.Conn) *Handle {
	return &Handle{
		userProvider:  *usecase.NewUser(conn),
		cacheProvider: *cache.NewCache(conn),
	}
}

func (handle *Handle) Login(c *gin.Context) {
	var user usecase.User
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	id, err := handle.userProvider.GetIDByLoginFromDB(c, user.Login, user.Password)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (handle *Handle) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(errors.Wrap(err, "GET /user/:id Atoi").Error())
		return
	}

	login, err := handle.userProvider.GetUserByIDFromDB(c, id)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": login})
}

func (handle *Handle) AddUser(c *gin.Context) {
	var user usecase.User

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	if err := handle.userProvider.AddingUserToDB(c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (handle *Handle) UpdateUser(c *gin.Context) {
	var user usecase.User

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "PUT /login ShouldBind").Error())
		return
	}

	if err := handle.userProvider.UpdateUserInDB(c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.Id), 10): "updated"})
}

func (handle *Handle) LoginCached(c *gin.Context) {
	var user usecase.User
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	id, err := handle.cacheProvider.GetIDByLoginFromDB(c, user.Login, user.Password)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (handle *Handle) GetUserByIDCached(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(errors.Wrap(err, "GET /user/:id Atoi").Error())
		return
	}

	login, err := handle.cacheProvider.GetUserByIDFromDB(c, id)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": login})
}
