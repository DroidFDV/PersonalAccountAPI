package handler

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Handle struct {
	userProvider  usecase.UserProvider
	workerManager *workers.Manager
}

func New(provider usecase.UserProvider, manager *workers.Manager) *Handle {
	return &Handle{
		userProvider:  provider,
		workerManager: manager,
	}
}

func (h *Handle) Login(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Handle.Login gin.ShouldBind").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	id, err := h.userProvider.GetIDByLogin(c, user)
	if err != nil {
		slog.Error(errors.Wrap(err, "Handle.Login userProvider.GetIDByLogin").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handle) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(errors.Wrap(err, "Handle.GetUserByID strconv.Atoi").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	login, err := h.userProvider.GetUserByID(c, models.UserRequest{ID: id})
	if err != nil {
		slog.Error(errors.Wrap(err, "Handle.GetUserByID userProvider.GetUserByID").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": login})
}

func (h *Handle) AddUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Handle.AddUser gin.ShouldBind").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	if err := h.userProvider.AddingUser(c, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error(errors.Wrap(err, "Handle.AddUser userProvider.AddingUser").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (h *Handle) UpdateUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Handle.UpdateUser gin.ShouldBind").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	if err := h.userProvider.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(errors.Wrap(err, "Handle.UpdateUser userProvider.UpdateUser").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.ID), 10): "updated"})
}

func (h *Handle) UploadFile(c *gin.Context) {
	file, err := c.FormFile("File")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		slog.Error(errors.Wrap(err, "Handle.UploadFile gin.FormFile").Error())
		return
	}

	h.workerManager.SetJob(h.userProvider.UploadFile(context.TODO(), file))
	// todo: сделать колбэк для неудачной загрузки

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    file.Filename,
	})
}
