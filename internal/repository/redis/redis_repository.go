package redis

import (
	"context"
	"time"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) models.TokenRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) SetTokens(accessKey string, accessExp time.Duration, refreshKey string, refreshExp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			var setErr error
			keys := map[string]time.Duration{
				accessKey:  accessExp,
				refreshKey: refreshExp,
			}

			for key, exp := range keys {
				setErr = pipe.Set(ctx, key, true, exp).Err()
				if setErr != nil {
					return setErr
				}
			}
			return nil
		})

		return err
	}, accessKey, refreshKey)
	return err
}

func (r *RedisRepository) Get(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := r.client.Get(ctx, key).Bool()

	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	return res, nil
}

func (r *RedisRepository) Count(pattern string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := r.client.Keys(ctx, pattern).Result()

	if err != nil {
		return 0, err
	}
	return len(res), nil
}
