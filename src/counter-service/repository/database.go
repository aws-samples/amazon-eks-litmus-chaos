package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.TODO()
)

const (
	countKey = "views"
)

type Database struct {
	Client *redis.Client
}

func NewDatabase(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &Database{
		Client: client,
	}, nil
}

func (d *Database) GetCount(ctx context.Context) (int64, error) {
	result, err := d.Client.Get(ctx, countKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

func (d *Database) IncrAndGetCount(ctx context.Context) (int64, error) {
	result, err := d.Client.Incr(ctx, countKey).Result()

	if err != nil {
		return 0, err
	}

	return result, nil
}
