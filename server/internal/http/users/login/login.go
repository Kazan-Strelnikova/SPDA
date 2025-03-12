package login

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (user.User, string, error)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func New(log *slog.Logger, service UserService, timeout time.Duration) func(c *gin.Context) {
	validate := validator.New()

	return func(c *gin.Context) {

		log.Info("login request")

		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Invalid request format", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Error("Validation failed", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		log = log.With(slog.String("email", req.Email))
		log.Info("attempt to log in user")

		user, token, err := service.Login(ctx, req.Email, req.Password)
		if err != nil {
			log.Error("Login failed", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		log.Info("login succeeded")

		c.SetCookie("token", token, 3600*24*365, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
