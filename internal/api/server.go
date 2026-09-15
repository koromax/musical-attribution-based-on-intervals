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
		c.Redirect(http.StatusFound, "/feed")
	})

	r.GET("/feed", h.GetFeed)
	r.GET("/feed/:id", h.GetFeed)
	r.GET("/draft", h.GetDraft)
	r.GET("/grid", h.GetGrid)

	r.Run(":8080")
}
