package enrollment

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UserEmail string
	EventId   uuid.UUID
}
