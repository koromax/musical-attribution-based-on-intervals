package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUserAPI(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	if login == "" || password == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.Repo.CreateUser(login, password)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"login": user.Login,
	})
}

func (h *Handler) LoginUserAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "аутентификация пройдена (нет)",
		"token":   "tokentokentoken",
	})
}

func (h *Handler) LogoutUserAPI(c *gin.Context) {
	c.Status(http.StatusOK)
}
