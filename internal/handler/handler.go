package handler

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Handle struct {
	userProvider  usecase.Provider
	workerManager *workers.Manager
}

func New(provider usecase.Provider, manager *workers.Manager) *Handle {
	return &Handle{
		userProvider:  provider,
		workerManager: manager,
	}
}

func (h *Handle) Login(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Handle.Login gin.ShouldBind:").Error())
		return
	}

	id, err := h.userProvider.GetIDByLoginFromDB(c, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		slog.Error(errors.Wrap(err, "Handle.Login userProvider.GetIDByLoginFromDB:").Error())
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
		slog.Error(errors.Wrap(err, "Handle.GetUserByID strconv.Atoi:").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	login, err := h.userProvider.GetUserByIDFromDB(c, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(errors.Wrap(err, "Handle.GetUserByID userProvider.GetUserByIDFromDB:").Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		slog.Error(errors.Wrap(err, "Handle.GetUserByID userProvider.GetUserByIDFromDB:").Error())
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
		slog.Error(errors.Wrap(err, "Handle.AddUser gin.ShouldBind:").Error())
		return
	}

	if err := h.userProvider.AddingUserToDB(c, user.ID, user.Login, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error(errors.Wrap(err, "Handle.AddUser userProvider.AddingUserToDB:").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (h *Handle) UpdateUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Handle.UpdateUser gin.ShouldBind:").Error())
		return
	}

	if err := h.userProvider.UpdateUserInDB(c, user.ID, user.Login, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(errors.Wrap(err, "Handle.UpdateUser userProvider.UpdateUserInDB:").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.ID), 10): "updated"})
}

func (h *Handle) UploadFile(c *gin.Context) {
	file, err := c.FormFile("File")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		slog.Error(errors.Wrap(err, "Handle.UploadFile gin.FormFile:").Error())
		return
	}

	h.userProvider.SetFile(file)
	h.workerManager.SetJob(h.userProvider.UploadFile)
	if err := h.workerManager.GetLog(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		slog.Error(errors.Wrap(err, "Handle.UploadFile workerManager.GetLog:").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    file.Filename,
	})
}
