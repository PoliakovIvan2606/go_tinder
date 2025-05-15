package apiserver

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tinder/internal/app/background"
	"tinder/internal/app/middleware"
	"tinder/internal/app/routes"
	"tinder/internal/app/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


type App struct {
	config *Config
	router *gin.Engine
	store *store.Store
	redis  *redis.Client
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
	a.configureRedis()

	worker := background.NewWorker(a.store, a.redis)

	if err := worker.LoadAllUserIDsToQueue(); err != nil {
		log.Fatal("failed to load user ids:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	worker.StartRecoWorker(ctx)

	// Настройка Gin
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
	routes.SetupRoutes(a.router, a.store, a.redis)

	// Запуск сервера в фоне
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- a.router.Run(addr)
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
		cancel()

		if err := a.redis.FlushDB(context.Background()).Err(); err != nil {
			log.Printf("Failed to flush Redis: %v", err)
		} else {
			log.Println("Redis flushed")
		}
		return nil

	case err := <-srvErr:
		log.Printf("Server error: %v", err)
		cancel()
		return err
	}
}

func (s *App) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err  
	}

	s.store = st

	return nil
}

func (s *App) configureRedis() {
	s.redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес Redis
		Password: "",               // без пароля, если не задан
		DB:       0,                // номер БД
	})
}