package get

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventService interface {
	GetEvent(ctx context.Context, eventID uuid.UUID) (event.Event, error)
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

		evt, err := service.GetEvent(ctx, eventID)
		if err != nil {
			log.Error("event not found", slog.String("event_id", eventIDStr), slog.String("error", err.Error()))
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}

		c.JSON(http.StatusOK, evt)
	}
}
