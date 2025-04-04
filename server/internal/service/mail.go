package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/letter"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage"
	"github.com/paulmach/orb"
	"gopkg.in/gomail.v2"
)

const DateFormat = "January 2, 2006 15:04"

type NominatimResponse struct {
	DisplayName string `json:"display_name"`
}

func (s *Service) SendEventChangesNotificationEmail(
	ctx context.Context,
	oldEvent event.Event,
	newEvent event.Event,
	name string,
	email string,
) error {
	const op = "service.SendNotificationEmail"
	log := s.log.With(
		slog.String("op", op),
		slog.String("event_id", oldEvent.ID.String()),
		slog.String("email", email),
	)

	log.Info("attempt to send notification email")

	changes := make([]letter.Change, 0, 2)

	if !oldEvent.Date.Equal(newEvent.Date) {
		changes = append(changes, letter.NewChange(
			"Date",
			oldEvent.Date.Format(DateFormat),
			newEvent.Date.Format(DateFormat),
		))
	}

	if !oldEvent.Location.Equal(newEvent.Location) {

		var oldLocation, newLocation string
		var err error

		if s.cache != nil {
			oldLocation, err = s.cache.GetLocation(ctx, oldEvent.Location)
			if err != nil {
				log.Error("error getting cache entry", slog.String("err", err.Error()))
				if !errors.Is(err, storage.ErrorLocationNotFound) {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}

				oldLocation, err = s.ReverseGeocode(oldEvent.Location)
				if err != nil {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}

				err = s.cache.SetLocation(ctx, oldEvent.Location, oldLocation)
				if err != nil {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}
			}

			newLocation, err = s.cache.GetLocation(ctx, newEvent.Location)
			if err != nil {
				log.Error("error getting cache entry", slog.String("err", err.Error()))
				if !errors.Is(err, storage.ErrorLocationNotFound) {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}

				newLocation, err = s.ReverseGeocode(newEvent.Location)
				if err != nil {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}

				err = s.cache.SetLocation(ctx, newEvent.Location, newLocation)
				if err != nil {
					return fmt.Errorf("op: %s, err: %v", op, err)
				}
			}
		} else {

			oldLocation, err = s.ReverseGeocode(oldEvent.Location)
			if err != nil {
				return fmt.Errorf("op: %s, err: %v", op, err)
			}

			newLocation, err = s.ReverseGeocode(newEvent.Location)
			if err != nil {
				return fmt.Errorf("op: %s, err: %v", op, err)
			}

		}

		changes = append(changes, letter.NewChange(
			"Location",
			oldLocation,
			newLocation,
		))
	}

	if len(changes) <= 0 {
		return nil
	}

	content := letter.NewUpdateNotification(name, oldEvent.Title, changes)

	return s.sendMessage("Changes in an upcoming event", content, email)
}

func (s *Service) SendEventNotificationEmail(
	ctx context.Context,
	name string,
	email string,
	evt event.Event,
) error {
	const op = "service.SendNotificationEmail"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("event_id", evt.ID.String()),
	)

	location, err := s.cache.GetLocation(ctx, evt.Location)
	if err != nil && !errors.Is(err, storage.ErrorLocationNotFound) {
		log.Error("error accessing cache", slog.String("err", err.Error()))
		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	if err != nil {
		location, err = s.ReverseGeocode(evt.Location)
		if err != nil {
			log.Error("error getting locaion from nominatim", slog.String("err", err.Error()))
			return fmt.Errorf("op: %s, err: %v", op, err)
		}

		err = s.cache.SetLocation(ctx, evt.Location, location)
		if err != nil {
			log.Error("error writing location to cache", slog.String("err", err.Error()))
			// return fmt.Errorf("op: %s, err: %v", op, err)
		}
	}

	content := letter.NewReminderNotification(
		name,
		evt.Title,
		location,
		evt.Date.Format(DateFormat),
	)

	return s.sendMessage("Upcoming event", content, email)
}

func (s *Service) ReverseGeocode(point orb.Point) (string, error) {
	const op = "service.ReverseGeocode"

	log := s.log.With(
		slog.String("op", op),
	)

	baseURL := "https://nominatim.openstreetmap.org/reverse"
	queryParams := url.Values{
		"lat":            {fmt.Sprintf("%f", point[0])},
		"lon":            {fmt.Sprintf("%f", point[1])},
		"format":         {"json"},
		"addressdetails": {"1"},
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	resp, err := http.Get(requestURL)
	if err != nil {
		log.Error("failed to fetch geolocation", slog.String("err", err.Error()))
		return "", fmt.Errorf("failed to fetch geolocation: %v", err)
	}
	defer resp.Body.Close()

	var result NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("failed to decode response", slog.String("err", err.Error()))
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if result.DisplayName == "" {
		return "Unknown location", nil
	}

	return result.DisplayName, nil
}

func (s *Service) sendMessage(subject, content, email string) error {
	const op = "service.sendMessage"

	m := gomail.NewMessage()

	m.SetHeader("From", s.smtp.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		s.smtp.Host,
		s.smtp.Port,
		s.smtp.Username,
		s.smtp.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	return nil
}
