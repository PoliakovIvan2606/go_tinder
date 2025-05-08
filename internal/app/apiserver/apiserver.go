package apiserver

import (
	"tinder/internal/app/middleware"
	"tinder/internal/app/routes"
	"tinder/internal/app/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


type App struct {
	config *Config
	router *gin.Engine
	store *store.Store
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run(addr string) error {
	if err := a.configureStore(); err != nil {
		return err
	}

	a.router = gin.Default()

	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	a.router.Use(middleware.AuthMiddleware())
	a.router.Static("/static", "./internal/static")
	a.router.LoadHTMLGlob("internal/templates/*")
	routes.SetupRoutes(a.router, a.store)

	return a.router.Run(addr)
}

func (s *App) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err  
	}

	s.store = st

	return nil
}
