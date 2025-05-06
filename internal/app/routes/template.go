package routes

import (
	"net/http"
	"tinder/internal/app/store"

	"github.com/gin-gonic/gin"
)

func SetupTemplateRoutes(group *gin.RouterGroup, templateHandler *TemplateHandler) {
	TemplateGroup := group.Group("/")
	{
		TemplateGroup.GET("/register", templateHandler.register)
		TemplateGroup.GET("/login", templateHandler.login)
		TemplateGroup.GET("/preferences/:id", templateHandler.preferences)
		TemplateGroup.GET("/user/:id", templateHandler.index)

	}
}

type TemplateHandler struct {
	st *store.Store
}

func NewTemplateHandler(st *store.Store) *TemplateHandler {
	return &TemplateHandler{st: st}
}

func(h *TemplateHandler) register(c *gin.Context) {
    c.HTML(http.StatusOK, "register.html", gin.H{})
}

func(h *TemplateHandler) preferences(c *gin.Context) {
    c.HTML(http.StatusOK, "preferences.html", gin.H{})
}

func(h *TemplateHandler) index(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{})
}

func(h *TemplateHandler) login(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{})
}