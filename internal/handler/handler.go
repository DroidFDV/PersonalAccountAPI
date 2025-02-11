package handler

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Handle struct {
	userProvider usecase.Provider
}

func New(provider usecase.Provider) *Handle {
	return &Handle{
		userProvider: provider,
	}
}

func (handle *Handle) Login(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "Login ShouldBind").Error())
		return
	}

	id, err := handle.userProvider.GetIDByLoginFromDB(c, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
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
		slog.Error(errors.Wrap(err, "GetUserByID Atoi").Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}

	login, err := handle.userProvider.GetUserByIDFromDB(c, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
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
	var user models.UserRequest

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "AddUser ShouldBind").Error())
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
	var user models.UserRequest

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "UpdateUser ShouldBind").Error())
		return
	}

	if err := handle.userProvider.UpdateUserInDB(c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.Id), 10): "updated"})
}

func (handle *Handle) UploadFile(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	url := c.PostForm("url")
	fileName := c.PostForm("fileName")

	if url == "" || fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо передать url и fileName"})
		return
	}

	// Создаем задачу
	job := workers.Job{
		ID:       time.Now().Nanosecond(),
		FileName: fileName,
		URL:      url,
	}

	// Отправляем в очередь на обработку
	jobs <- job
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Файл %s добавлен в очередь", fileName)})
	// file, err := c.FormFile("File")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
	// 	return
	// }

	// dstPath, err := createSaveFilePath(userId, file.Filename)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
	// 	return
	// }
	// if err := c.SaveUploadedFile(file, dstPath); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "File uploaded successfully",
	// 	"file":    file.Filename,
	// })
}

func createSaveFilePath(dstDir string, fileName string) (string, error) {
	userDir := filepath.Join(models.UploadsDir, dstDir)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "Failed to create uploads directory: ")
	}

	dstPath := filepath.Join(userDir, fileName)
	return dstPath, nil
}
