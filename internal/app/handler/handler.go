package handler

import (
	"interval_attribution/internal/app/ds"
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

const (
	DefaultImage  = "http://localhost:9000/music/default_cover.jpg"
	DefaultVideo  = "http://localhost:9000/music/default_video.mp4"
	DefaultUserID = 1
)

type ComposerView struct {
	ds.Composer
	LikesCount int64
}

func prepareMedia(image, video string) (string, string) {
	if image == "" {
		image = DefaultImage
	}
	if video == "" {
		video = DefaultVideo
	}
	return image, video
}

func (h *Handler) GetComposersFeed(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Некорректный ID")
			return
		}
	}

	next := c.Query("next") == "true"

	composer, err := h.Repo.GetComposerFeedItem(id, next)
	if err != nil {
		c.String(http.StatusNotFound, "Опубликованная карточка не найдена")
		return
	}

	composer.ImageURL, composer.VideoURL = prepareMedia(composer.ImageURL, composer.VideoURL)
	likes := h.Repo.GetLikesCount(composer.ID)

	view := ComposerView{
		Composer:   composer,
		LikesCount: likes,
	}

	c.HTML(http.StatusOK, "composers-feed.html", gin.H{
		"composer": view,
	})
}

func (h *Handler) GetComposerDraft(c *gin.Context) {
	composer, err := h.Repo.GetComposerDraft(DefaultUserID)
	if err != nil {
		composer = ds.Composer{}
	}

	composer.ImageURL, composer.VideoURL = prepareMedia(composer.ImageURL, composer.VideoURL)
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

	var composerViews []ComposerView
	for _, comp := range composers {
		comp.ImageURL, comp.VideoURL = prepareMedia(comp.ImageURL, comp.VideoURL)
		likes := h.Repo.GetLikesCount(comp.ID)

		composerViews = append(composerViews, ComposerView{
			Composer:   comp,
			LikesCount: likes,
		})
	}

	c.HTML(http.StatusOK, "composers-grid.html", gin.H{
		"composers": composerViews,
		"query":     queryStr,
		"max_query": maxQueryStr,
	})
}

func (h *Handler) CreateComposerDraft(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		name = "Новый композитор"
	}

	_, err := h.Repo.CreateDraftComposer(DefaultUserID, name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/composer-draft")
}

func (h *Handler) PublishComposer(c *gin.Context) {
	idStr := c.PostForm("id")
	id, _ := strconv.Atoi(idStr)

	description := c.PostForm("description")
	freq1, _ := strconv.ParseFloat(c.PostForm("freq1"), 64)
	freq2, _ := strconv.ParseFloat(c.PostForm("freq2"), 64)

	err := h.Repo.PublishComposer(uint(id), description, freq1, freq2)
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка при публикации")
		return
	}

	c.Redirect(http.StatusFound, "/composers-grid")
}

func (h *Handler) DeleteComposer(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Некорректный ID")
		return
	}

	err = h.Repo.DeleteComposerSQL(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при удалении: "+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/composers-grid")
}
