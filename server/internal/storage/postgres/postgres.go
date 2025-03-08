package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxInterface interface {
    Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Close()
}

type Storage struct {
    conn    PgxInterface
}

func New(ctx context.Context, connStr string) *Storage {
	const op = "storage.postgres.New"

	defer ctx.Done()

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