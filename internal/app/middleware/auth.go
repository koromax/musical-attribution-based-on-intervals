package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"

	"interval_attribution/internal/app/ds"
	appRedis "interval_attribution/internal/app/redis"
)

func AuthMiddleware(jwtSecret string, rClient *appRedis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверка блеклиста в Redis
		err := rClient.CheckJWTInBlacklist(c.Request.Context(), tokenStr)
		if err == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is revoked (logged out)"})
			return
		}

		if !errors.Is(err, redis.Nil) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("is_moderator", claims.IsModerator)

		c.Next()
	}
}
