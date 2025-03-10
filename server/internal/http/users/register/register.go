package register

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
	Register(ctx context.Context, usr user.User) (string, error)
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	LastName string `json:"last_name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func New(log *slog.Logger, service UserService, timeout time.Duration) func(c *gin.Context) {	
	validate := validator.New()

	return func(c *gin.Context) {
		var req RegisterRequest
		
		log.Info("registration request")

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
		
		usr := user.User{
			Name:     req.Name,
			LastName: req.LastName,
			Email:    req.Email,
			Password: req.Password, 
		}

		log = log.With(slog.String("email", req.Email))
		log.Info("attempt to register user")
		
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		token, err := service.Register(ctx, usr)
		if err != nil {
			log.Error("Registration failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
			return
		}
		
		log.Info("registration succeeded")

		usr.Password = ""
		
		c.SetCookie("token", token, 3600, "/", "", false, true)
		
		c.JSON(http.StatusCreated, gin.H{
			"user": usr,
		})
	}
}
