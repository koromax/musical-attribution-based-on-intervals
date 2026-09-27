package api

import (
	"interval_attribution/internal/app/handler"
	"interval_attribution/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	repo := repository.NewRepository()
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

	r.Run(":8080")
}
