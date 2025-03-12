package create

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paulmach/orb"
)

type EventService interface {
	CreateEvent(ctx context.Context, evt event.Event) (event.Event, error)
}

type CreateEventRequest struct {
	Title        string `json:"title" binding:"required" validate:"required"`
	Type         int    `json:"type" binding:"required" validate:"required"`
	Date         string `json:"date" binding:"required" validate:"required"`
	TotalSeats   int    `json:"total_seats" binding:"required,min=1" validate:"required,min=1"`
	CreatorEmail string `json:"creator_email" binding:"required,email" validate:"required,email"`
	Location     struct {
		Latitude  string `json:"latitude" binding:"required" validate:"required"`
		Longitude string `json:"longitude" binding:"required" validate:"required"`
	} `json:"location" binding:"required" validate:"required"`
	Description *string `json:"description,omitempty"`
}

func New(log *slog.Logger, service EventService, timeout time.Duration) func(c *gin.Context) {
	validate := validator.New()
	return func(c *gin.Context) {
		var req CreateEventRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(400, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		lat, err := strconv.ParseFloat(req.Location.Latitude, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid latitude"})
			return
		}

		long, err := strconv.ParseFloat(req.Location.Longitude, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid longitude"})
			return
		}

		location := orb.Point{lat, long}

		eventDate, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format, expected RFC3339"})
			return
		}

		evt := event.Event{
			Title:          req.Title,
			Type:           req.Type,
			Date:           eventDate,
			TotalSeats:     req.TotalSeats,
			AvailableSeats: req.TotalSeats,
			CreatorEmail:   req.CreatorEmail,
			Location:       location,
			Description:    req.Description,
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		createdEvent, err := service.CreateEvent(ctx, evt)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error creating event", "details": err.Error()})
			return
		}

		c.JSON(200, createdEvent)
	}
}
