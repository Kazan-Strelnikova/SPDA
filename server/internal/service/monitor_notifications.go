package service

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
)

type UserNotificationRequest struct {
	Email string
	Evt   event.Event
}

const monitoringCount = 4

func (s *Service) Monitor(done <-chan os.Signal) {
	const op = "service.Monitor"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("starting monitoring")

	var lastChecked = time.Now()

	userNotificationQueue := make(chan UserNotificationRequest, monitoringCount)
	eventQueue := make(chan event.Event, monitoringCount)

	defer close(eventQueue)
	defer close(userNotificationQueue)

	for range monitoringCount {
		go s.notifyOfEventWorker(eventQueue, userNotificationQueue)
		go s.notifyUserWorker(userNotificationQueue)
	}

	for {
		select {
		case <-done:
			return
		default:
			diff := time.Since(lastChecked)
			var sleepTime time.Duration

			if diff < 20*time.Minute {
				sleepTime = 405 * 3 * time.Second
			} else {
				sleepTime = 395 * 3 * time.Second
			}

			lastChecked = time.Now()

			before := lastChecked.Add((24*60 + 10) * time.Minute)
			after := lastChecked.Add((24*60 - 10) * time.Minute)

			eventsToNotifyOf, err := s.GetAllEvents(
				context.Background(),
				nil,
				nil,
				nil,
				nil,
				&before,
				&after,
				nil,
				nil,
				nil,
			)

			log.Debug(
				"some debug message",
				slog.Time("now", time.Now()),
				slog.Any("events", eventsToNotifyOf),
				slog.Time("before", before),
				slog.Time("after", after),
				slog.Time("lastChecked", lastChecked),
			)

			if err != nil {
				log.Error(
					"error getting events to notify of",
					slog.String("err", err.Error()),
				)
				time.Sleep(sleepTime)
				continue
			}

			for _, event := range eventsToNotifyOf {
				eventQueue <- event
			}

			time.Sleep(sleepTime)
		}

	}
}

func (s *Service) notifyOfEventWorker(eventQueue <-chan event.Event, userNotificationQueue chan<- UserNotificationRequest) {
	const op = "service.NotifyOfEvent"

	log := s.log.With(
		slog.String("op", op),
	)

	for evt := range eventQueue {
		log.Info(
			"notifying of event",
			slog.String("event_id", evt.ID.String()),
		)
		usersToNotifyOfEvent, err := s.evtRepo.GetAllSubscriptions(
			context.Background(),
			evt.ID,
		)

		if err != nil {
			log.Error("error retrieving enrolled users", slog.String("err", err.Error()))
			continue
		}

		for _, user := range usersToNotifyOfEvent {
			userNotificationQueue <- UserNotificationRequest{
				Email: user.UserEmail,
				Evt:   evt,
			}
		}
	}

}

func (s *Service) notifyUserWorker(taskQueue <-chan UserNotificationRequest) {
	const op = "service.notifyUserWorker"

	log := s.log.With(
		slog.String("op", op),
	)

	for task := range taskQueue {
		log.Info(
			"notifying user of event",
			slog.String("emain", task.Email),
			slog.String("event_id", task.Evt.ID.String()),
		)
		usr, err := s.usrRepo.GetUser(context.Background(), task.Email)
		if err != nil {
			log.Error("error getting user",
				slog.String("email", task.Email),
				slog.String("err", err.Error()),
			)

			continue
		}

		err = s.SendEventNotificationEmail(
			context.Background(),
			usr.Name,
			usr.Email,
			task.Evt,
		)

		if err != nil {
			log.Error("error notifying user of event",
				slog.String("email", task.Email),
				slog.String("err", err.Error()),
			)
		}
	}
}
