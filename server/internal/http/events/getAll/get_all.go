package getall

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
)

type EventService interface {
	GetAllEvents(
		ctx context.Context,
		limit, offset *int,
		eventType *int,
		creatorEmail *string,
		before, after *time.Time,
		location *orb.Point,
		radius *float64,
		visitorEmail *string,
	) ([]event.Event, error)
}

func New(log *slog.Logger, service EventService, timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		var (
			limit, offset *int
			eventType     *int
			creatorEmail  *string
			before, after *time.Time
			location      *orb.Point
			radius        *float64
			visitorEmail  *string
		)

		// Parse query params
		if val, exists := c.GetQuery("limit"); exists {
			if num, err := strconv.Atoi(val); err == nil {
				limit = &num
			}
		}

		if val, exists := c.GetQuery("offset"); exists {
			if num, err := strconv.Atoi(val); err == nil {
				offset = &num
			}
		}

		if val, exists := c.GetQuery("type"); exists {
			if num, err := strconv.Atoi(val); err == nil {
				eventType = &num
			}
		}

		if val, exists := c.GetQuery("creator_email"); exists {
			creatorEmail = &val
		}

		if val, exists := c.GetQuery("before"); exists {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				before = &t
			}
		}

		if val, exists := c.GetQuery("after"); exists {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				after = &t
			}
		}

		if latStr, latExists := c.GetQuery("lat"); latExists {
			if lonStr, lonExists := c.GetQuery("lon"); lonExists {
				if lat, err1 := strconv.ParseFloat(latStr, 64); err1 == nil {
					if lon, err2 := strconv.ParseFloat(lonStr, 64); err2 == nil {
						if val, exists := c.GetQuery("radius"); exists {
							if num, err := strconv.ParseFloat(val, 64); err == nil {
								loc := orb.Point{lon, lat}
								location = &loc
								radius = &num
							}
						}
					}
				}
			}
		}

		if val, exists := c.GetQuery("visitor_email"); exists {
			visitorEmail = &val
		}

		events, err := service.GetAllEvents(ctx, limit, offset, eventType, creatorEmail, before, after, location, radius, visitorEmail)
		if err != nil {
			log.Error("error retrieving events", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
			return
		}

		c.JSON(http.StatusOK, events)
	}
}
