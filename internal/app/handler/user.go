package handler

import "github.com/gin-gonic/gin"

func getUserID(c *gin.Context) uint {
	val, exists := c.Get("user_id")
	if !exists {
		return 0 // 0 для неавторизованных
	}
	if id, ok := val.(uint); ok {
		return id
	}
	return 0
}
