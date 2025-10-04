package handler

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.Login gin.ShouldBind", slog.Any("error", err))
		return
	}

	id, err := h.userProvider.GetIDByLogin(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.Login userProvider.GetIDByLogin", slog.Any("error", err))
		return
	}
	if id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handle.Login authentication failed", slog.String("reason", "invalid credentials"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handle) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.GetUserByID strconv.Atoi", slog.Any("error", err))
		return
	}

	login, err := h.userProvider.GetUserByID(c, models.UserRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.GetUserByID userProvider.GetUserByID", slog.Any("error", err))
		return
	}
	if login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handle.GetUserByID authorization failed", slog.String("reason", "user not found or access denied"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": login})
}

func (h *Handle) AddUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.AddUser gin.ShouldBind", slog.Any("error", err))
		return
	}

	if err := h.userProvider.AddingUser(c, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error("Handle.AddUser userProvider.AddingUser", slog.Any("error", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (h *Handle) UpdateUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handle.UpdateUser gin.ShouldBind", slog.Any("error", err))
		return
	}

	if err := h.userProvider.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handle.UpdateUser userProvider.UpdateUser", slog.Any("error", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.ID), 10): "updated"})
}

func (h *Handle) UploadFile(c *gin.Context) {
	file, err := c.FormFile("File")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		slog.Error("Handle.UploadFile gin.FormFile", slog.Any("error", err))
		return
	}

	h.workerManager.SetJob(h.userProvider.UploadFile(c, file))

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    file.Filename,
	})
}
