package handler

import (
	"interval_attribution/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repo: r}
}

func (h *Handler) GetComposersFeed(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	var err error
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
	}

	next := c.Query("next") == "true"

	composer, err := h.Repo.GetComposerFeedItem(id, next)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.HTML(http.StatusOK, "composers-feed.html", gin.H{"composer": composer})
}

func (h *Handler) GetComposerDraft(c *gin.Context) {
	composer, err := h.Repo.GetComposerDraft()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.HTML(http.StatusOK, "composer-draft.html", gin.H{"composer": composer})
}

func (h *Handler) GetComposersGrid(c *gin.Context) {
	queryStr := c.Query("query")
	maxQueryStr := c.Query("max_query")

	var minFreq, maxFreq *float64

	if val, err := strconv.ParseFloat(queryStr, 64); err == nil && queryStr != "" {
		minFreq = &val
	}

	if val, err := strconv.ParseFloat(maxQueryStr, 64); err == nil && maxQueryStr != "" {
		maxFreq = &val
	}

	composers := h.Repo.GetPublishedComposers(minFreq, maxFreq)

	c.HTML(http.StatusOK, "composers-grid.html", gin.H{
		"composers": composers,
		"query":     queryStr,
		"max_query": maxQueryStr,
	})
}
