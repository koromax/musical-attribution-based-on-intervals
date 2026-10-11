package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"interval_attribution/internal/app/ds"
)

func (h *Handler) GetComposersAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

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

func (h *Handler) GetComposerFeedAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

	currentID, _ := strconv.Atoi(c.Query("id"))
	next := c.Query("next") == "true"

	composer, err := h.Repo.GetComposerFeedItem(currentID, next, currentUserID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, composer)
}

func (h *Handler) GetComposerDraftAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

	composer, err := h.Repo.GetComposerDraft(currentUserID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, composer)
}

func (h *Handler) CreateComposerAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID
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

func (h *Handler) PublishComposerAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

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

func (h *Handler) DeleteComposerAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

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

func (h *Handler) LikeComposerAPI(c *gin.Context) {
	currentUserID := ds.GetCurrentUserSingleton().UserID

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
