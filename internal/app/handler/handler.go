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

func (h *Handler) GetFeed(c *gin.Context) {
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

	item, err := h.Repo.GetFeedItem(id, next)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.HTML(http.StatusOK, "feed.html", gin.H{"item": item})
}

func (h *Handler) GetDraft(c *gin.Context) {
	draft, err := h.Repo.GetDraft()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.HTML(http.StatusOK, "draft.html", gin.H{"item": draft})
}

func (h *Handler) GetGrid(c *gin.Context) {
	queryStr := c.Query("query")
	maxQueryStr := c.Query("max_query")

	var minFreq, maxFreq *float64

	if val, err := strconv.ParseFloat(queryStr, 64); err == nil && queryStr != "" {
		minFreq = &val
	}

	if val, err := strconv.ParseFloat(maxQueryStr, 64); err == nil && maxQueryStr != "" {
		maxFreq = &val
	}

	items := h.Repo.GetPublishedGrid(minFreq, maxFreq)

	c.HTML(http.StatusOK, "grid.html", gin.H{
		"items":     items,
		"query":     queryStr,
		"max_query": maxQueryStr,
	})
}
