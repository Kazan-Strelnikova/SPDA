package ping

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	}
}
