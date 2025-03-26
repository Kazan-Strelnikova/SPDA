package service

import (
	"errors"
	"log/slog"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/config"
)

type Service struct {
	log     *slog.Logger
	usrRepo UserRepository
	evtRepo EventRepository
	tknScrt string
	smtp    config.SMTPConfig
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
) *Service {
	return &Service{
		log:     log,
		usrRepo: usrRepo,
		evtRepo: evtRepo,
		tknScrt: tokenSecret,
		smtp:    smtp,
	}
}
