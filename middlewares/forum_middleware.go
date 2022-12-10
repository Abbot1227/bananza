package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		postId := c.Params.ByName("id")
		// Validate post id
		if len(postId) != 27 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Post id is not valid"})
		}
		c.Next()
	}
}
