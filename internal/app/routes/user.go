package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	// "net/http"
	"strconv"
	"tinder/internal/app/models"
	"tinder/internal/app/store"
	"tinder/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupUserRoutes(group *gin.RouterGroup, userHandler *UserHandler) {
	userGroup := group.Group("/user")
	{
		userGroup.GET("/", userHandler.GetUser)
		userGroup.POST("/add", userHandler.CreateUser)
		userGroup.POST("/login", userHandler.CheckUser)
		userGroup.GET("/preferences", userHandler.GetPreferencesUser)
	}
}

type UserHandler struct {
	redis *redis.Client
	st *store.Store
}

func NewUserHandler(st *store.Store, redis *redis.Client) *UserHandler {
	return &UserHandler{
		st: st, 
		redis: redis,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	idStr, ok := idRaw.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID type"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	a, err := h.st.User().UserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, gin.H{
		"user": a,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.UserCreate
	if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": "Неверный JSON"})
		return
	}

	var id int
	id, err := h.st.User().Create(&user)
	if err != nil {
		c.JSON(422, gin.H{
			"error": err.Error(),
			"example_request": gin.H{
				"name":        "John Doe",
				"email":       "john@example.com",
				"password":    "your_password",
				"age":         25,
				"description": "I'm from Russia",
				"city":        "Moscow",
				"latitude":    55.7558,
				"longitude":   37.6176,
			},
		})
		return
	}
	c.JSON(201, gin.H{"user_id": id}) 
}

func (h *UserHandler) CheckUser(c *gin.Context) {
	var userCheck models.UserCheck
	if err := c.ShouldBindJSON(&userCheck); err != nil {
		c.JSON(400, gin.H{"error": "Неверный JSON"})
		return
	}

	id, passwordHash, err := h.st.User().IdAndPaswordByEmail(userCheck.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	if id == "" || !models.CheckPasswordHash(passwordHash, userCheck.Password) {
		c.JSON(401, gin.H{"error": "Неверный email или пароль"})
		return
	}

	accessToken, err := jwt.GenerateAccessToken(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось создать access токен"})
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось создать refresh токен"})
		return
	}

	c.SetCookie("access", accessToken, 900, "/", "", true, true)      // 900 = 15 мин
	c.SetCookie("refresh", refreshToken, 604800, "/", "", true, true) // 7 дней

	c.JSON(200, gin.H{"message": "Добро пожаловать"})
}

func (h *UserHandler) GetPreferencesUser(c *gin.Context) {
	idRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	idStrAuth, ok := idRaw.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID type"})
		return
	}

	id, err := strconv.Atoi(idStrAuth)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()
	key := fmt.Sprintf("user:%d:preferences", id)

	idStr, err := h.redis.LPop(ctx, key).Result()
	if err == redis.Nil {
		// Очередь пуста
		c.JSON(200, gin.H{"message": "No more preferences"})
		return
	}
	if err != nil {
		log.Println("LPOP error:", err)
		c.JSON(500, gin.H{"error": "Redis error"})
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid user ID:", idStr)
		c.JSON(400, gin.H{"error": "Invalid user ID in queue"})
		return
	}

	a, err := h.st.User().UserById(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, gin.H{
		"user": a,
	})
}