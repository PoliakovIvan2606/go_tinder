package routes

import (
	"database/sql"
	"errors"
	"strconv"
	"tinder/internal/app/models"
	"tinder/internal/app/store"

	"github.com/gin-gonic/gin"
)

func SetPreferencesRoutes(group *gin.RouterGroup, preferencesHandler *PreferencesHandler) {
	preferencesGroup := group.Group("/preferences")
	{
		preferencesGroup.GET("/:id", preferencesHandler.GetPreferences)
		preferencesGroup.POST("/add", preferencesHandler.CreatePreferences)
	}
}

type PreferencesHandler struct {
	st *store.Store
}

func NewPreferencesHandler(st *store.Store) *PreferencesHandler {
	return &PreferencesHandler{st: st}
}

func (h *PreferencesHandler) GetPreferences(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	a, err := h.st.Preferences().GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "Preferences not found"})
			return
		}
		c.JSON(400, gin.H{"error": "Failed to get preferences"})
		return
	}

	c.JSON(200, gin.H{
		"preferences": a,
	})
}

func (h *PreferencesHandler) CreatePreferences(c *gin.Context) {
	var preferences models.Preferences
	if err := c.ShouldBindJSON(&preferences); err != nil {
        c.JSON(400, gin.H{"error": "Неверный JSON"})
		return
	}

	var id int
	id, err := h.st.Preferences().Create(&preferences)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
			"example_request": gin.H{
				"user_id":		1,
				"gender":      	"Мужчина|Женщина",
				"age_from":		18,
				"age_to":		25,
				"radius":		100,
			},
		})
		return
	}
	c.JSON(201, gin.H{"preferences_id": id}) 
}