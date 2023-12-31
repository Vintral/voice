package utils

import (
	"github.com/gin-gonic/gin"
)

func GetEmailFromContext(c *gin.Context) string {
	test, _ := c.Get("email")
	if val, ok := test.(string); ok {
		return val
	}

	return ""
}
