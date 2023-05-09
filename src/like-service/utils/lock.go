package utils

import (
	"like-service/common"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type Database struct {
	Client  *goredislib.Client
	Redsync *redsync.Redsync
}

func NewDatabase(address string) (*Database, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	if err := client.Ping(common.Ctx).Err(); err != nil {
		return nil, err
	}

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	return &Database{
		Client:  client,
		Redsync: rs,
	}, nil
}
