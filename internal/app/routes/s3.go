package routes

import (
	"context"
	"net/http"
	"strconv"
	"tinder/internal/app/models"
	"tinder/internal/app/store"
	"tinder/pkg/s3"

	"github.com/gin-gonic/gin"
)

func SetupS3Routes(group *gin.RouterGroup, S3Handler *S3Handler) {
	S3Group := group.Group("/")
	{
		S3Group.POST("/upload/:id", S3Handler.uploadMultiple)
	}
}

type S3Handler struct {
	st *store.Store
	S3 *s3.S3
}

func NewS3Handler(st *store.Store) *S3Handler {
	return &S3Handler{
		st: st,
		S3: s3.NewS3(),
	}
}

func (h *S3Handler) uploadMultiple(c *gin.Context) {
    id := c.Param("id")

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения файлов"})
        return
    }

    files := form.File["files"] 
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Файлы не получены"})
        return
    }

    ctx := context.TODO()
    urls := make([]string, 0, len(files))

    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть файл"})
            return
        }

        fileHeader.Filename = id + "/" + fileHeader.Filename

        url, err := h.S3.Upload(ctx, fileHeader, file)
        file.Close()

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки в S3", "details": err.Error()})
            return
        }
        urls = append(urls, url)
    }

	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный user_id"})
		return
	}
	
	user_photo := models.Photo{
		User_id: userID,
		Photos:  urls,
	}
	
	if err := h.st.Photo().Create(&user_photo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки в БД"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Фото успешно сохранены"})
}
