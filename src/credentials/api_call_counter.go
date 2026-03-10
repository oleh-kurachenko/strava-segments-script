package credentials

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const counterReadKey = "stravaapi:counter:read"
const counterReadLimit = 100

type APICallCounter struct {
	redisClient *redis.Client
}

func NewAPICallCounter(client *redis.Client) (counter *APICallCounter) {
	counter = new(APICallCounter)
	counter.redisClient = client

	return counter
}

func (counter *APICallCounter) getValue() (counterVal int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	counterTTL, err := counter.redisClient.TTL(ctx, counterReadKey).Result()
	if err != nil {
		return 0, err
	}

	if counterTTL == -2 {
		return 0, nil
	}

	return counter.redisClient.Get(ctx, counterReadKey).Int()
}

func (counter *APICallCounter) IsFine() (bool, error) {
	counterVal, err := counter.getValue()
	if err != nil {
		return false, err
	}

	return counterVal < counterReadLimit, nil
}

func (counter *APICallCounter) Increment() error {
	counterVal, err := counter.getValue()
	if err != nil {
		return err
	}

	counterVal++
	tTL := (time.Now().Truncate(15 * time.Minute).Add(15 * time.Minute)).Sub(
		time.Now())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return counter.redisClient.Set(ctx, counterReadKey, counterVal, tTL).Err()
}
