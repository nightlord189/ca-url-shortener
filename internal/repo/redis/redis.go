package redis

import (
	"context"
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

const expiration = 1 * time.Hour

type Repo struct {
	client *redis.Client
}

func New(cfg config.RedisConfig) (*Repo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.Database, // use default DB
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Repo{client: client}, nil
}

func (r *Repo) PutLink(ctx context.Context, shortURL, originalURL string) error {
	return r.client.Set(ctx, shortURL, originalURL, expiration).Err()
}

func (r *Repo) GetLink(ctx context.Context, shortURL string) (string, error) {
	return r.client.Get(ctx, shortURL).Result()
}
