package routes

import (
	"tinder/internal/app/store"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(router *gin.Engine, st *store.Store, redis *redis.Client) {
	user_heandler := NewUserHandler(st, redis)
	preferences_heandler := NewPreferencesHandler(st)
	api := router.Group("/api")
	{
		SetupUserRoutes(api, user_heandler)  // /api/users
		SetPreferencesRoutes(api, preferences_heandler)  // /api/preferences
	}
	template_heandler := NewTemplateHandler(st)
	template := router.Group("/")
	{
		SetupTemplateRoutes(template, template_heandler)  
	}
	s3_heandler := NewS3Handler(st)
	s3 := router.Group("/")
	{
		SetupS3Routes(s3, s3_heandler)
	}
}