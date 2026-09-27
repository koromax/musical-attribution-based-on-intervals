package api

import (
	"log"
	"net/http"

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

	repo := repository.NewRepository(db)
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/composers-feed")
	})

	r.GET("/composers-feed", h.GetComposersFeed)
	r.GET("/composers-feed/:id", h.GetComposersFeed)
	r.GET("/composer-draft", h.GetComposerDraft)
	r.GET("/composers-grid", h.GetComposersGrid)

	r.POST("/composer/create", h.CreateComposerDraft)
	r.POST("/composer/publish", h.PublishComposer)
	r.POST("/composer/delete", h.DeleteComposer)

	r.Run(":8080")
}
