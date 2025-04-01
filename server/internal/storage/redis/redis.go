package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/event"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(
	addr string,
) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:                  addr,
		DB:                    0,
		ContextTimeoutEnabled: true,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetEvent(ctx context.Context, id uuid.UUID) (event.Event, error) {
	const op = "storage.redis.GetEvent"
	key := "ev" + id.String()

	cmd := c.client.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err().Error() == redis.Nil.Error() {
			return event.Event{}, storage.ErrorEventNotFound
		}
		return event.Event{}, cmd.Err()
	}

	cmdb, err := cmd.Bytes()
	if err != nil {
		return event.Event{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	b := bytes.NewReader(cmdb)

	var res event.Event

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return event.Event{}, fmt.Errorf("op: %s, err: %v", op, err)
	}

	return res, nil
}

func (c *Client) SetEvent(ctx context.Context, evt event.Event) error {
	const op = "storage.redis.SetEvent"

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(evt); err != nil {
		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	return c.client.Set(ctx, "ev"+evt.ID.String(), b.Bytes(), time.Until(evt.Date)).Err()
}

func (c *Client) GetLocation(ctx context.Context, point orb.Point) (string, error) {
	const op = "storage.redis.GetLocation"

	key := "pt" + fmt.Sprintf("%f%f", point.Lat(), point.Lon())

	cmd := c.client.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err().Error() == redis.Nil.Error() {
			return "", storage.ErrorLocationNotFound
		}
		return "", cmd.Err()
	}

	cmdb, err := cmd.Bytes()
	if err != nil {
		return "", fmt.Errorf("op: %s, err: %v", op, err)
	}

	b := bytes.NewReader(cmdb)

	var res string

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return "", fmt.Errorf("op: %s, err: %v", op, err)
	}

	return res, nil
}

func (c *Client) SetLocation(ctx context.Context, point orb.Point, content string) error {
	const op = "storage.redis.SetLocation"

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(content); err != nil {
		return fmt.Errorf("op: %s, err: %v", op, err)
	}

	return c.client.Set(ctx, "pt"+fmt.Sprintf("%f%f", point.Lat(), point.Lon()), b.Bytes(), 0).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
