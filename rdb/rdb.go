package rdb

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var rdb *redis.Client
var ctx context.Context

func init() {
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     consts.REDIS_ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Debug().Msgf("start rdb %s", consts.REDIS_ADDR)
}

func Instance() *redis.Client {
	return rdb
}

func PublishEmail(email *mailer.RedisQueueEmail) error {
	payload, err := json.Marshal(email)

	if err != nil {
		return err
	}

	return Publish(mailer.REDIS_EMAIL_CHANNEL, payload)
}

func Publish(channel string, data []byte) error {
	//log.Debug().Msgf("send %v", data)
	return rdb.Publish(ctx, channel, data).Err()
}
