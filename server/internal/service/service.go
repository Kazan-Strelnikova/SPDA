package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/config"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Cache interface {
	GetEvent(ctx context.Context, id uuid.UUID) (event.Event, error)
	SetEvent(ctx context.Context, evt event.Event) error
	GetLocation(ctx context.Context, point orb.Point) (string, error)
	SetLocation(ctx context.Context, point orb.Point, content string) error
}

type Service struct {
	log     *slog.Logger
	usrRepo UserRepository
	evtRepo EventRepository
	tknScrt string
	smtp    config.SMTPConfig
	cache 	Cache

}

var (
	ErrInvalidToken   = errors.New("invalid user token")
	ErrUserNotFound   = errors.New("user with given email does not exist")
	ErrEventNotFound  = errors.New("event with this id does not exist")
	ErrNotEnoughSeats = errors.New("not enough seats")
)

func New(
	log *slog.Logger,
	usrRepo UserRepository,
	evtRepo EventRepository,
	tokenSecret string,
	smtp config.SMTPConfig,
	cache Cache,
) *Service {
	return &Service{
		log:     log,
		usrRepo: usrRepo,
		evtRepo: evtRepo,
		tknScrt: tokenSecret,
		smtp:    smtp,
		cache: 	 cache,
	}
}
