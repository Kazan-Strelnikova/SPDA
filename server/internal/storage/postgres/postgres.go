package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/enrollment"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/orb"
)

type PgxInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type Storage struct {
	conn PgxInterface
}

func New(ctx context.Context, connStr string) *Storage {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) GetUser(ctx context.Context, email string) (user.User, error) {
	const op = "storage.postgres.GetUser"

	var usr user.User

	query := `
	SELECT name, last_name, email, password
	FROM users
	WHERE email=$1
	`

	err := s.conn.QueryRow(ctx, query, email).Scan(
		&usr.Name,
		&usr.LastName,
		&usr.Email,
		&usr.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, storage.ErrorNoUser
		}

		return user.User{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	return usr, nil
}

func (s *Storage) InsertUser(ctx context.Context, usr user.User) error {
	const op = "storage.postgres.InsertUser"

	query := `
	INSERT INTO users(name, last_name, email, password)
	VALUES ($1, $2, $3, $4)
	`

	_, err := s.conn.Exec(ctx, query,
		usr.Name,
		usr.LastName,
		usr.Email,
		usr.Password,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return storage.ErrorUserExists
			}
		}

		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	return nil
}

func (s *Storage) InsertEvent(ctx context.Context, evt *event.Event) error {
	const op = "storage.postgres.InsertEvent"

	query := `
	INSERT INTO events (title, type, date, total_seats, available_seats, creator_email, location, description)
	VALUES ($1, $2, $3, $4, $5, $6, ST_GeomFromText($7, 4326), $8)
	RETURNING id
	`

	err := s.conn.QueryRow(ctx, query,
		evt.Title,
		evt.Type,
		evt.Date,
		evt.TotalSeats,
		evt.AvailableSeats,
		evt.CreatorEmail,
		fmt.Sprintf("POINT(%f %f)", evt.Location.Lon(), evt.Location.Lat()),
		evt.Description,
	).Scan(&evt.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return storage.ErrorNoUser
			}
		}

		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	const op = "storage.postgres.DeleteEvent"

	query := `DELETE FROM events WHERE id = $1`

	cmdTag, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return storage.ErrorEventNotFound
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (event.Event, error) {
	const op = "storage.postgres.GetEvent"

	var ev event.Event
	var locationWKT string

	query := `
	SELECT id, title, type, date, total_seats, available_seats, creator_email, 
	       ST_AsText(location), description
	FROM events
	WHERE id = $1
	`

	err := s.conn.QueryRow(ctx, query, id).Scan(
		&ev.ID,
		&ev.Title,
		&ev.Type,
		&ev.Date,
		&ev.TotalSeats,
		&ev.AvailableSeats,
		&ev.CreatorEmail,
		&locationWKT,
		&ev.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return event.Event{}, storage.ErrorEventNotFound
		}
		return event.Event{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	ev.Location, err = parseWKT(locationWKT)
	if err != nil {
		return event.Event{}, fmt.Errorf("op: %s, err: invalid location format: %v", op, err)
	}

	return ev, nil
}

func (s *Storage) GetAllEvents(ctx context.Context, limit, offset *int, eventType *int, creatorEmail *string, before, after *time.Time, location *orb.Point, radius *float64, visitorEmail *string) ([]event.Event, error) {
	const op = "storage.postgres.GetAllEvents"

	var events = []event.Event{}

	query := `
	SELECT e.id, e.title, e.type, e.date, e.total_seats, e.available_seats, e.creator_email, ST_AsText(e.location), e.description
	FROM events e
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if eventType != nil {
		conditions = append(conditions, fmt.Sprintf("e.type = $%d", argCount))
		args = append(args, *eventType)
		argCount++
	}

	if creatorEmail != nil {
		conditions = append(conditions, fmt.Sprintf("e.creator_email = $%d", argCount))
		args = append(args, *creatorEmail)
		argCount++
	}

	if before != nil {
		conditions = append(conditions, fmt.Sprintf("e.date < $%d", argCount))
		args = append(args, *before)
		argCount++
	}

	if after != nil {
		conditions = append(conditions, fmt.Sprintf("e.date > $%d", argCount))
		args = append(args, *after)
		argCount++
	}

	if location != nil && radius != nil {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(e.location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326), $%d)", argCount, argCount+1, argCount+2))
		args = append(args, location.Lon(), location.Lat(), *radius)
		argCount += 3
	}

	if visitorEmail != nil {
		conditions = append(conditions, fmt.Sprintf(`
		e.id IN (
			SELECT event_id FROM enrollments WHERE user_email = $%d
		)`, argCount))
		args = append(args, *visitorEmail)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limit != nil && offset != nil {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, limit, offset)
	}

	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %v", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var ev event.Event
		var locationWKT string
		err := rows.Scan(
			&ev.ID,
			&ev.Title,
			&ev.Type,
			&ev.Date,
			&ev.TotalSeats,
			&ev.AvailableSeats,
			&ev.CreatorEmail,
			&locationWKT,
			&ev.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("op: %s, err: %v", op, err)
		}
		ev.Location, err = parseWKT(locationWKT)
		if err != nil {
			return nil, fmt.Errorf("op: %s, err: invalid location format: %v", op, err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("op: %s, err: %v", op, err)
	}

	return events, nil
}

func (s *Storage) SubscribeToEvent(ctx context.Context, eventId uuid.UUID, email string) (uuid.UUID, error) {
	const op = "storage.postgres.SubscribeToEvent"

	var id uuid.UUID

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	query := `
	INSERT INTO enrollments(user_email, event_id)
	VALUES ($1, $2)
	RETURNING id
	`

	err = tx.QueryRow(ctx, query, email, eventId).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	query = `
	UPDATE events
	SET available_seats = available_seats - 1
	WHERE id = $1
	`

	_, err = tx.Exec(ctx, query, eventId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: failed to commit transaction: %v", op, err)
	}

	return id, nil
}

func (s *Storage) GetEventSubscription(ctx context.Context, enrollmentEventId uuid.UUID, email string) (enrollment.Enrollment, error) {
	const op = "storage.postgres.GetEventSubscription"

	var enrlmnt enrollment.Enrollment

	query := `
	SELECT id, created_at, user_email, event_id
	FROM enrollments
	WHERE user_email = $1 AND event_id = $2
	`

	var enrlmntId, eventId string

	err := s.conn.QueryRow(ctx, query, email, enrollmentEventId).Scan(
		&enrlmntId,
		&enrlmnt.CreatedAt,
		&enrlmnt.UserEmail,
		&eventId,
	)

	if err != nil {
		return enrollment.Enrollment{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	enrlmnt.Id, err = uuid.Parse(enrlmntId)
	if err != nil {
		return enrollment.Enrollment{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	enrlmnt.EventId, err = uuid.Parse(eventId)
	if err != nil {
		return enrollment.Enrollment{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	return enrlmnt, nil
}

func (s *Storage) UnsubscribeFromEvent(ctx context.Context, enrollmentId uuid.UUID) error {
	const op = "storage.postgres.UnsibscribeFromEvent"

	query := `
	DELETE FROM enrollments
	WHERE id=$1
	`

	_, err := s.conn.Exec(ctx, query, enrollmentId)
	if err != nil {
		return fmt.Errorf("op: %s, err: %v", op, err)
	}
	
	return nil
}

func parseWKT(wkt string) (orb.Point, error) {
	// fmt.Println(wkt)
	parts := strings.Fields(strings.TrimPrefix(strings.TrimSuffix(wkt, ")"), "POINT("))
	if len(parts) != 2 {
		return orb.Point{}, fmt.Errorf("invalid WKT format: %s", wkt)
	}

	var lon, lat float64
	_, err := fmt.Sscanf(parts[0], "%f", &lon)
	if err != nil {
		return orb.Point{}, err
	}

	_, err = fmt.Sscanf(parts[1], "%f", &lat)
	if err != nil {
		return orb.Point{}, err
	}

	return orb.Point{lon, lat}, nil
}
