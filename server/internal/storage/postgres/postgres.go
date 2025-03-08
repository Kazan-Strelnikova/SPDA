package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxInterface interface {
    Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type Storage struct {
    conn    PgxInterface
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