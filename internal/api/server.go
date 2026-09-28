package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"interval_attribution/internal/app/dsn"
	"interval_attribution/internal/app/handler"
	"interval_attribution/internal/app/repository"
)

func StartServer() {
	_ = godotenv.Load()

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

	h := handler.NewHandler(repo)
	r := gin.Default()

	api := r.Group("/api")
	{
		// Домен композиторов (услуг)
		api.GET("/composers", h.GetComposersAPI)
		api.GET("/composer/feed", h.GetComposerFeedAPI)
		api.GET("/composer/draft", h.GetComposerDraftAPI)
		api.POST("/composer", h.CreateComposerAPI)
		api.PUT("/composer/:id/publish", h.PublishComposerAPI)
		api.DELETE("/composer/:id", h.DeleteComposerAPI)
		api.POST("/composer/:id/like", h.LikeComposerAPI)

		// Домен пользователя
		api.POST("/user/register", h.RegisterUserAPI)
		api.POST("/user/login", h.LoginUserAPI)
		api.POST("/user/logout", h.LogoutUserAPI)
	}

	r.Run(":8080")
}
