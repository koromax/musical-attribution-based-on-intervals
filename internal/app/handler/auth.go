package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"interval_attribution/internal/app/config"
	"interval_attribution/internal/app/ds"
	"interval_attribution/internal/app/redis"
	"interval_attribution/internal/app/repository"
)

type AuthHandler struct {
	repo        *repository.Repository
	redisClient *redis.Client
	cfg         *config.Config
}

func NewAuthHandler(repo *repository.Repository, rClient *redis.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		repo:        repo,
		redisClient: rClient,
		cfg:         cfg,
	}
}

type registerReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register — регистрация нового пользователя
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body registerReq true "Данные для регистрации"
// @Success 200 {object} map[string]interface{}
// @Router /user/register [post]
func (h *AuthHandler) Register(gCtx *gin.Context) {
	var req registerReq
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.repo.CreateUser(req.Login, req.Password)
	if err != nil {
		gCtx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login — вход и выдача JWT
// @Summary Авторизация (вход)
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginReq true "Учетные данные"
// @Success 200 {object} map[string]interface{}
// @Router /user/login [post]
func (h *AuthHandler) Login(gCtx *gin.Context) {
	var req loginReq
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByLogin(req.Login)
	if err != nil || user.Password != req.Password {
		gCtx.JSON(http.StatusForbidden, gin.H{"error": "invalid login or password"})
		return
	}

	claims := &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(h.cfg.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID:      user.ID,
		IsModerator: user.IsModerator,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.cfg.JWT.Token))
	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{
		"access_token": tokenStr,
		"token_type":   "Bearer",
		"expires_in":   int64(h.cfg.JWT.ExpiresIn.Seconds()),
	})
}

// Logout — выход (добавление токена в Redis blacklist)
// @Summary Выход (деавторизация)
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /user/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсим токен, чтобы узнать оставшееся время жизни (TTL) для блеклиста в Redis
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &ds.JWTClaims{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token claims"})
		return
	}

	// Вычисляем оставшееся время жизни токена
	expTime := time.Unix(claims.ExpiresAt, 0)
	ttl := time.Until(expTime)
	if ttl <= 0 {
		ttl = time.Minute // минимальный TTL, если токен уже просрочен
	}

	// Записываем токен в блеклист Redis
	err = h.redisClient.WriteJWTToBlacklist(c.Request.Context(), tokenStr, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
