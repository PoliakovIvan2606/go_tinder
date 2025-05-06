package routes

import (
	"database/sql"
	"errors"
	"strconv"
	"tinder/internal/app/models"
	"tinder/internal/app/store"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(group *gin.RouterGroup, userHandler *UserHandler) {
	userGroup := group.Group("/user")
	{
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.POST("/add", userHandler.CreateUser)
		userGroup.POST("/login", userHandler.CheckUser)
	}
}

type UserHandler struct {
	st *store.Store
}

func NewUserHandler(st *store.Store) *UserHandler {
	return &UserHandler{st: st}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")

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
				"age":         25,
				"description": "I'm from Russia",
				"city":        "Moscow",
				"coordinates": "(55.7558, 37.6176)",
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

	user, err := h.st.User().UserByEmail(userCheck.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	if user == nil {
		c.JSON(401, gin.H{"error": "Пользователь не найден"})
		return
	}
	if !user.CheckPasswordHash(userCheck.Password) {
		c.JSON(401, gin.H{"error": "Неправильный пароль"})
		return
	}

	c.JSON(200, gin.H{"message": "Добро пожаловать"})
}
