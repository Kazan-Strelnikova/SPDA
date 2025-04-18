package create

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentService interface {
	SubscribeToEvent(ctx context.Context, eventId uuid.UUID, email string) (uuid.UUID, error)
}

func New(log *slog.Logger, service AppointmentService, timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Info("create appointment request")

		emailAny, ok := c.Get("email")
		if !ok {
			log.Error("email does not exist on the context")
			c.JSON(http.StatusBadRequest, gin.H{"error": "email does not exist on the context"})
			return
		}

		email, ok := emailAny.(string)
		if !ok {
			log.Error("email is not of proper format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is not of proper format"})
			return
		}

		eventId, err := uuid.Parse(c.Param("event_id"))
		if err != nil {
			log.Error("Invalid event id", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		id, err := service.SubscribeToEvent(ctx, eventId, email)
		if err != nil {
			log.Error("Could not subscribe to event")
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not subscribe to event"})
			return
		}

		log.Info("create appointment request succeeded")

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}
