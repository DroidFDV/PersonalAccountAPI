package handler

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/uploading"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
)

type Handler struct {
	userProvider   usecase.UserProvider
	workerManager  *workers.Manager
	uploadProvider uploading.UploadingProvider
}

func New(provider usecase.UserProvider, manager *workers.Manager, uploadProvider uploading.UploadingProvider) *Handler {
	return &Handler{
		userProvider:   provider,
		workerManager:  manager,
		uploadProvider: uploadProvider,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.Login gin.ShouldBind", slog.Any("error", err))
		return
	}

	userResponce, err := h.userProvider.GetIDByLogin(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.Login userProvider.GetIDByLogin", slog.Any("error", err))
		return
	}
	if userResponce.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handler.Login authentication failed", slog.String("reason", "invalid credentials"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": userResponce.ID})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.GetUserByID strconv.Atoi", slog.Any("error", err))
		return
	}

	userResponce, err := h.userProvider.GetUserByID(c, models.UserRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.GetUserByID userProvider.GetUserByID", slog.Any("error", err))
		return
	}
	if userResponce.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handler.GetUserByID authorization failed", slog.String("reason", "user not found or access denied"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponce.Login})
}

func (h *Handler) AddUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.AddUser gin.ShouldBind", slog.Any("error", err))
		return
	}

	if err := h.userProvider.AddUser(c, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error("Handler.AddUser userProvider.AddingUser", slog.Any("error", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.UpdateUser gin.ShouldBind", slog.Any("error", err))
		return
	}

	if err := h.userProvider.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error("Handler.UpdateUser userProvider.UpdateUser", slog.Any("error", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.ID), 10): "updated"})
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		slog.Error("Handler.UploadFile gin.FormFile", slog.Any("error", err))
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		slog.Error("Handler.UploadFile gin.Param", slog.String("reason", "no parameters received"))
		return
	}

	//NOTE: вопрос о сохранении файла
	reqCtx := c.Request.Context()
	h.workerManager.SetJob(func(ctx context.Context) error {
		if err := h.uploadProvider.Upload(reqCtx, id, file); err != nil {
			return errors.Wrap(err, "SetJob uploadProvider.Upload")
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    file.Filename,
	})
}
