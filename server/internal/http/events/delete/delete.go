package delete

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventService interface {
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
}

func New(log *slog.Logger, service EventService, timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		eventIDStr := c.Param("event_id")
		eventID, err := uuid.Parse(eventIDStr)
		if err != nil {
			log.Error("invalid event ID", slog.String("event_id", eventIDStr), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
			return
		}

		email, ok := c.Get("email")
		if !ok {
			log.Error("email does not exist on the context")
			c.JSON(http.StatusBadRequest, gin.H{"error": "email does not exist on the context"})
		}

		err = service.DeleteEvent(context.WithValue(ctx, "email", email), eventID)
		if err != nil {
			log.Error("event not found", slog.String("event_id", eventIDStr), slog.String("error", err.Error()))
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
