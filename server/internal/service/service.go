package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/enrollment"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	log     *slog.Logger
	usrRepo UserRepository
	evtRepo EventRepository
	tknScrt string
}

type UserRepository interface {
	GetUser(ctx context.Context, email string) (user.User, error)
	InsertUser(ctx context.Context, usr user.User) error
}

type EventRepository interface {
	InsertEvent(ctx context.Context, evt *event.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (event.Event, error)
	GetAllEvents(ctx context.Context, limit, offset *int, eventType *int, creatorEmail *string, before, after *time.Time, location *orb.Point, radius *float64, visitorEmail *string) ([]event.Event, error)
	SubscribeToEvent(ctx context.Context, eventId uuid.UUID, email string) (uuid.UUID, error)
	UnsubscribeFromEvent(ctx context.Context, enrollmentId uuid.UUID) error
	GetEventSubscription(ctx context.Context, enrollmentId uuid.UUID, email string) (enrollment.Enrollment, error)
}

var (
	ErrInvalidToken = errors.New("invalid user token")
)

func New(
	log *slog.Logger,
	usrRepo UserRepository,
	evtRepo EventRepository,
	tokenSecret string,
) *Service {
	return &Service{
		log:     log,
		usrRepo: usrRepo,
		evtRepo: evtRepo,
		tknScrt: tokenSecret,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (user.User, string, error) {
	const op = "service.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	usr, err := s.usrRepo.GetUser(ctx, email)
	if err != nil {
		log.Error("error getting user", slog.String("error", err.Error()))
		return user.User{}, "", fmt.Errorf("invalid email or password")
	}

	if !checkPasswordHash(password, usr.Password) {
		log.Error("incorrect password")
		return user.User{}, "", fmt.Errorf("invalid email or password")
	}

	usr.Password = ""

	token, err := generateTokens(s.tknScrt, usr)
	if err != nil {
		log.Error("error creating token", slog.String("error", err.Error()))
		return user.User{}, "", fmt.Errorf("error creating token")
	}

	return usr, token, nil
}

func (s *Service) LoginByToken(ctx context.Context, token string) (user.User, error) {
	const op = "service.LoginByToken"

	log := s.log.With(
		slog.String("op", op),
	)

	secret := []byte(s.tknScrt)

	data, err := jwt.ParseWithClaims(token, &user.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok && token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return secret, nil
		}
		return nil, ErrInvalidToken
	})

	if err != nil {
		log.Error("invalid jwt token", slog.String("err", err.Error()))
		return user.User{}, ErrInvalidToken
	}

	if claims, ok := data.Claims.(*user.UserClaims); ok && data.Valid {
		log.Info("token validated successfully", slog.Any("email", claims.Payload.Email))
		return claims.Payload, nil
	}

	log.Error("invalid jwt token")
	return user.User{}, ErrInvalidToken
}

func (s *Service) Register(ctx context.Context, usr user.User) (string, error) {
	const op = "service.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", usr.Email),
	)

	var err error

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		log.Error("error hashing password", slog.String("error", err.Error()))
		return "", fmt.Errorf("invalid password")
	}

	err = s.usrRepo.InsertUser(ctx, usr)
	if err != nil {
		log.Error("error creating user", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating user")
	}

	usr.Password = ""

	token, err := generateTokens(s.tknScrt, usr)
	if err != nil {
		log.Error("error creating user", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating user")
	}

	return token, nil
}

func (s *Service) CreateEvent(ctx context.Context, evt event.Event) (event.Event, error) {
	const op = "service.CreateEvent"

	log := s.log.With(
		slog.String("op", op),
		slog.String("creator_email", evt.CreatorEmail),
	)

	_, err := s.usrRepo.GetUser(ctx, evt.CreatorEmail)
	if err != nil {
		log.Error("creator not found", slog.String("error", err.Error()))
		return event.Event{}, fmt.Errorf("creator not found")
	}

	err = s.evtRepo.InsertEvent(ctx, &evt)
	if err != nil {
		log.Error("error inserting event", slog.String("error", err.Error()))
		return event.Event{}, fmt.Errorf("error creating event")
	}

	return evt, nil
}

func (s *Service) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	const op = "service.DeleteEvent"

	log := s.log.With(
		slog.String("op", op),
		slog.String("event_id", eventID.String()),
	)

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		log.Error("email not found in context")
		return fmt.Errorf("email not found in context")
	}

	evt, err := s.evtRepo.GetEvent(ctx, eventID)
	if err != nil {
		log.Error("error getting event", slog.String("error", err.Error()))
		return fmt.Errorf("event not found")
	}

	if evt.CreatorEmail != email {
		log.Error("user is not the creator of the event")
		return fmt.Errorf("user is not the creator of this event")
	}

	err = s.evtRepo.DeleteEvent(ctx, eventID)
	if err != nil {
		log.Error("error deleting event", slog.String("error", err.Error()))
		return fmt.Errorf("error deleting event")
	}

	return nil
}

func (s *Service) GetEvent(ctx context.Context, eventID uuid.UUID) (event.Event, error) {
	const op = "service.GetEvent"

	log := s.log.With(
		slog.String("op", op),
		slog.String("event_id", eventID.String()),
	)

	evt, err := s.evtRepo.GetEvent(ctx, eventID)
	if err != nil {
		log.Error("error getting event", slog.String("error", err.Error()))
		return event.Event{}, fmt.Errorf("event not found")
	}

	return evt, nil
}

func (s *Service) GetAllEvents(
	ctx context.Context,
	limit, offset *int,
	eventType *int,
	creatorEmail *string,
	before, after *time.Time,
	location *orb.Point,
	radius *float64,
	visitorEmail *string,
) ([]event.Event, error) {
	const op = "service.GetAllEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.evtRepo.GetAllEvents(
		ctx,
		limit,
		offset,
		eventType,
		creatorEmail,
		before,
		after,
		location,
		radius,
		visitorEmail,
	)
	if err != nil {
		log.Error("error retrieving events", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error retrieving events")
	}

	return events, nil
}

func (s *Service) SubscribeToEvent(ctx context.Context, eventId uuid.UUID, email string) (uuid.UUID, error) {
	const op = "service.SubscribeToEvent"

	var id uuid.UUID

	log := s.log.With(
		slog.String("event_id", eventId.String()),
		slog.String("email", email),
		slog.String("op", op),
	)

	_, err := s.usrRepo.GetUser(ctx, email)
	if err != nil {
		log.Error("error retrieving user", slog.String("err", err.Error()))
		return uuid.UUID{}, fmt.Errorf("error retrieving user")
	}

	evt, err := s.evtRepo.GetEvent(ctx, eventId)
	if err != nil {
		log.Error("error retrieving event", slog.String("err", err.Error()))
		return uuid.UUID{}, fmt.Errorf("error retrieving event")
	}

	if evt.AvailableSeats < 1 {
		log.Error("not enough free spaces on the event")
		return uuid.UUID{}, fmt.Errorf("not enough free spaces on the event")
	}

	id, err = s.evtRepo.SubscribeToEvent(ctx, eventId, email)
	if err != nil {
		log.Error("error making an appointment", slog.String("err", err.Error()))
		return uuid.UUID{}, fmt.Errorf("error making an appointment")
	}

	return id, nil
}

func (s *Service) UnsibscribeFromEvent(ctx context.Context, eventId uuid.UUID, email string) error {
	const op = "service.UnsibscribeFromEvent"

	log := s.log.With(
		slog.String("email", email),
		slog.String("event_id", eventId.String()),
		slog.String("op", op),
	)

	enrlmnt, err := s.evtRepo.GetEventSubscription(ctx, eventId, email)
	if err != nil {
		log.Error("error accessing event enrollment", slog.String("err", err.Error()))
		return fmt.Errorf("error accessing event enrollment")
	}

	err = s.evtRepo.UnsubscribeFromEvent(ctx, enrlmnt.Id)
	if err != nil {
		log.Error("error deleting enrollment", slog.Any("enrollment", enrlmnt))
		return fmt.Errorf("error deleting enrollment")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateTokens(accessSecret string, usr user.User) (string, error) {

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(),
		"payload": usr,
	}).SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
