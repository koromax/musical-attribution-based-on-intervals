package api

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "interval_attribution/docs"
	"interval_attribution/internal/app/config"
	"interval_attribution/internal/app/dsn"
	"interval_attribution/internal/app/handler"
	"interval_attribution/internal/app/middleware"
	appRedis "interval_attribution/internal/app/redis"
	"interval_attribution/internal/app/repository"
)

func StartServer() {
	_ = godotenv.Load()
	ctx := context.Background()

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	repo, err := repository.NewRepository(&repository.RepositorySettings{
		DB:              db,
		MinioEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		MinioBucketName: os.Getenv("MINIO_BUCKET_NAME"),
	})
	if err != nil {
		log.Fatalf("Ошибка MinIO/Repository: %v", err)
	}

	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	rClient, err := appRedis.New(
		ctx,
		os.Getenv("REDIS_HOST"),
		redisPort,
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_USER"),
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	defer rClient.Close()

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("Критическая ошибка: JWT_SECRET_KEY не задан в переменных окружения")
	}

	h := handler.NewHandler(repo)
	authHandler := handler.NewAuthHandler(repo, rClient, &config.Config{
		JWT: config.JWTConfig{
			Token:     jwtSecret,
			ExpiresIn: 24 * time.Hour,
		},
	})
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/composers", h.GetComposersAPI)
		api.GET("/composer/feed", h.GetComposerFeedAPI)

		api.POST("/user/register", authHandler.Register)
		api.POST("/user/login", authHandler.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret, rClient))
		{
			protected.GET("/composer/draft", h.GetComposerDraftAPI)
			protected.POST("/composer", h.CreateComposerAPI)
			protected.PUT("/composer/:id/publish", h.PublishComposerAPI)
			protected.DELETE("/composer/:id", h.DeleteComposerAPI)
			protected.POST("/composer/:id/like", h.LikeComposerAPI)

			protected.POST("/user/logout", authHandler.Logout)
		}
	}

	r.Run(":8080")
}
