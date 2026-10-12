package api

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}
	if cfg.JWT.Token == "" {
		log.Fatal("Критическая ошибка: JWT_SECRET_KEY не задан в переменных окружения")
	}

	ctx := context.Background()

	db, err := gorm.Open(postgres.Open(dsn.FromConfig(cfg.Database)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	repo, err := repository.NewRepository(&repository.RepositorySettings{
		DB:              db,
		MinioEndpoint:   cfg.Minio.Endpoint,
		MinioAccessKey:  cfg.Minio.AccessKey,
		MinioSecretKey:  cfg.Minio.SecretKey,
		MinioBucketName: cfg.Minio.Bucket,
	})
	if err != nil {
		log.Fatalf("Ошибка MinIO/Repository: %v", err)
	}

	rClient, err := appRedis.New(
		ctx,
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.User,
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	defer rClient.Close()

	h := handler.NewHandler(repo)
	authHandler := handler.NewAuthHandler(repo, rClient, cfg)
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/composers", h.GetComposersAPI)
		api.GET("/composer/feed", h.GetComposerFeedAPI)

		api.POST("/user/register", authHandler.Register)
		api.POST("/user/login", authHandler.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Token, rClient))
		{
			protected.GET("/composer/draft", h.GetComposerDraftAPI)
			protected.POST("/composer", h.CreateComposerAPI)
			protected.PUT("/composer/:id/publish", h.PublishComposerAPI)
			protected.DELETE("/composer/:id", h.DeleteComposerAPI)
			protected.POST("/composer/:id/like", h.LikeComposerAPI)

			protected.POST("/user/logout", authHandler.Logout)
		}
	}

	if err := r.Run(net.JoinHostPort(cfg.ServiceHost, strconv.Itoa(cfg.ServicePort))); err != nil {
		log.Fatalf("Не удалось запустить HTTP-сервер: %v", err)
	}
}
