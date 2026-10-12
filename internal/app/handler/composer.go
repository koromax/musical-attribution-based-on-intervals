package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetComposersAPI — получение списка опубликованных композиторов
// @Summary Получение списка опубликованных композиторов
// @Tags composers
// @Produce json
// @Param min_freq query number false "Минимальная частота"
// @Param max_freq query number false "Максимальная частота"
// @Success 200 {object} []ds.Composer
// @Router /composers [get]
func (h *Handler) GetComposersAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	var minFreq, maxFreq *float64
	if val := c.Query("min_freq"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			minFreq = &f
		}
	}
	if val := c.Query("max_freq"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			maxFreq = &f
		}
	}

	composers, err := h.Repo.GetPublishedComposers(minFreq, maxFreq, currentUserID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, composers)
}

// GetComposerFeedAPI — получение ленты композиторов
// @Summary Получение ленты композиторов
// @Tags composers
// @Produce json
// @Param id query int false "Текущий ID"
// @Param next query boolean false "Запрос далее"
// @Success 200 {object} ds.Composer
// @Router /composer/feed [get]
func (h *Handler) GetComposerFeedAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	currentID, _ := strconv.Atoi(c.Query("id"))
	next := c.Query("next") == "true"

	composer, err := h.Repo.GetComposerFeedItem(currentID, next, currentUserID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, composer)
}

// GetComposerDraftAPI — получение черновика композитора
// @Summary Получение черновика композитора
// @Tags composers
// @Produce json
// @Success 200 {object} ds.Composer
// @Router /composer/draft [get]
func (h *Handler) GetComposerDraftAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	composer, err := h.Repo.GetComposerDraft(currentUserID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, composer)
}

// CreateComposerAPI — создание нового черновика композитора с медиафайлами
// @Summary Создание черновика композитора
// @Tags composers
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Имя композитора"
// @Param pic formData file false "Изображение"
// @Param video formData file false "Видео"
// @Success 201 {object} ds.Composer
// @Router /composer [post]
func (h *Handler) CreateComposerAPI(c *gin.Context) {
	currentUserID := getUserID(c)
	name := c.PostForm("name")

	if name == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	composer, err := h.Repo.CreateDraftComposer(currentUserID, name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var imageURL, videoURL string

	if imageHeader, err := c.FormFile("pic"); err == nil {
		if url, err := h.Repo.UploadToMinIO(imageHeader, "img"); err == nil {
			imageURL = url
		}
	}

	if videoHeader, err := c.FormFile("video"); err == nil {
		if url, err := h.Repo.UploadToMinIO(videoHeader, "vid"); err == nil {
			videoURL = url
		}
	}

	if imageURL != "" || videoURL != "" {
		_ = h.Repo.UpdateComposerMedia(composer.ID, imageURL, videoURL)
		composer.ImageURL = imageURL
		composer.VideoURL = videoURL
	}

	c.JSON(http.StatusCreated, composer)
}

// PublishComposerAPI — публикация черновика композитора
// @Summary Публикация черновика композитора
// @Tags composers
// @Accept multipart/form-data
// @Security BearerAuth
// @Param id path int true "ID композитора"
// @Param description formData string false "Описание"
// @Param freq1 formData number false "Частота 1"
// @Param freq2 formData number false "Частота 2"
// @Success 200
// @Router /composer/{id}/publish [put]
func (h *Handler) PublishComposerAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	description := c.PostForm("description")
	freq1, _ := strconv.ParseFloat(c.PostForm("freq1"), 64)
	freq2, _ := strconv.ParseFloat(c.PostForm("freq2"), 64)

	err = h.Repo.PublishComposer(uint(id), currentUserID, description, freq1, freq2)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteComposerAPI — мягкое удаление композитора
// @Summary Удаление композитора
// @Tags composers
// @Security BearerAuth
// @Param id path int true "ID композитора"
// @Success 200
// @Router /composer/{id} [delete]
func (h *Handler) DeleteComposerAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteComposerSoft(uint(id), currentUserID)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
}

// LikeComposerAPI — установка или удаление лайка для композитора
// @Summary Лайк/Дизлайк композитора
// @Tags composers
// @Accept multipart/form-data
// @Security BearerAuth
// @Param id path int true "ID композитора"
// @Param like formData int true "Значение лайка"
// @Success 200
// @Router /composer/{id}/like [post]
func (h *Handler) LikeComposerAPI(c *gin.Context) {
	currentUserID := getUserID(c)

	idParam := c.Param("id")
	composerID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	likeVal, _ := strconv.Atoi(c.PostForm("like"))

	err = h.Repo.ToggleLike(currentUserID, uint(composerID), likeVal)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
