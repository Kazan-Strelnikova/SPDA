package tokenlogin

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	LoginByToken(ctx context.Context, token string) (user.User, error)
}

func New(log *slog.Logger, service UserService, timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		
		log.Info("attempt of cookie login")

		token, err := c.Cookie("token")
		if err != nil {
			log.Error("Missing authentication token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		usr, err := service.LoginByToken(ctx, token)
		if err != nil {
			log.Error("Invalid authentication token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Info("cookie login successful", slog.String("email", usr.Email))
		
		c.JSON(http.StatusOK, gin.H{"user": usr})
	}
}
