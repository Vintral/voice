package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

func GetEmailFromContext(c *gin.Context) string {
	test, _ := c.Get("email")
	if val, ok := test.(string); ok {
		return val
	}

	return ""
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
