package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Event struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	Type              int       `json:"type"`
	Date              time.Time `json:"date"`
	TotalSeats        int       `json:"total_seats"`
	AvailableSeats    int       `json:"available_seats"`
	CreatorEmail      string    `json:"creator_email"`
	Location          orb.Point `json:"location"`
	HasUnlimitedSeats string    `json:"has_unlimited_seats"`
	Description       *string   `json:"description,omitempty"`
}
