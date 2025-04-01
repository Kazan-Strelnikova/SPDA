package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/enrollment"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type EventRepository interface {
	InsertEvent(ctx context.Context, evt *event.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (event.Event, error)
	GetAllEvents(ctx context.Context, limit, offset *int, eventType *int, creatorEmail *string, before, after *time.Time, location *orb.Point, radius *float64, visitorEmail *string) ([]event.Event, error)
	SubscribeToEvent(ctx context.Context, eventId uuid.UUID, email string) (uuid.UUID, error)
	UnsubscribeFromEvent(ctx context.Context, enrollmentId uuid.UUID) error
	GetEventSubscription(ctx context.Context, enrollmentId uuid.UUID, email string) (enrollment.Enrollment, error)
	UpdateEvent(ctx context.Context, evt event.Event) error
	GetAllSubscriptions(ctx context.Context, eventId uuid.UUID) ([]enrollment.Enrollment, error)
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
		return event.Event{}, ErrUserNotFound
	}

	err = s.evtRepo.InsertEvent(ctx, &evt)
	if err != nil {
		log.Error("error inserting event", slog.String("error", err.Error()))
		return event.Event{}, fmt.Errorf("error creating event")
	}

	return evt, nil
}

func (s *Service) UpdateEvent(ctx context.Context, evt event.Event) error {
	const op = "service.UpdateEvent"

	log := s.log.With(
		slog.String("op", op),
		slog.String("creator_email", evt.CreatorEmail),
		slog.String("event_id", evt.ID.String()),
	)

	_, err := s.usrRepo.GetUser(ctx, evt.CreatorEmail)
	if err != nil {
		log.Error("creator not found", slog.String("error", err.Error()))
		return ErrUserNotFound
	}

	evt.AvailableSeats = evt.TotalSeats

	oldEvent, err := s.evtRepo.GetEvent(ctx, evt.ID)
	if err != nil {
		return ErrEventNotFound
	}

	enrolledUsers, err := s.evtRepo.GetAllSubscriptions(ctx, evt.ID)
	if err != nil {
		log.Error("could not retrieve enrollments", slog.String("err", err.Error()))
		return fmt.Errorf("could not retrieve enrollments")
	}

	if evt.HasUnlimitedSeats != "true" {

		if len(enrolledUsers) > evt.TotalSeats {
			return ErrNotEnoughSeats
		}

		evt.AvailableSeats -= len(enrolledUsers)
	}

	err = s.evtRepo.UpdateEvent(ctx, evt)
	if err != nil {
		log.Error("could not update event", slog.String("error", err.Error()))
		switch {
		default:
			return fmt.Errorf("error updating event")
		case errors.Is(err, storage.ErrorEventNotFound):
			return ErrEventNotFound
		}
	}

	for _, enrollment := range enrolledUsers {
		enrollment := enrollment
		user, err := s.usrRepo.GetUser(ctx, enrollment.UserEmail)
		if err != nil {
			return fmt.Errorf("op: %s, err: %v", op, err)
		}

		var newCtx context.Context
		var cancel context.CancelFunc

		deadline, ok := ctx.Deadline()
		if ok {
			newCtx, cancel = context.WithDeadline(context.WithoutCancel(ctx), deadline)
			_ = cancel
		} else {
			newCtx = context.WithoutCancel(ctx)
		}

		go s.SendEventChangesNotificationEmail(
			newCtx,
			oldEvent,
			evt,
			user.Name,
			user.Email,
		)
	}

	return nil
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
