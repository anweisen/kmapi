package cache

import (
  "context"
  "encoding/json"
  "errors"
  "github.com/redis/go-redis/v9"
  "os"
  "time"
)

type RedisCache struct {
  RedisClient *redis.Client
  Context     context.Context
}

func NewRedisCache() Cache {
  client := redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_URL"),
  })

  return RedisCache{
    client,
    context.Background(),
  }
}

func (c RedisCache) SetJson(key string, value any) error {
  data, err := json.Marshal(value)
  if err != nil {
    return err
  }

  return c.RedisClient.Set(c.Context, key, data, 24*time.Hour).Err()
}

func (c RedisCache) GetJson(key string, dest any) error {
  data, err := c.RedisClient.Get(c.Context, key).Bytes()
  if errors.Is(err, redis.Nil) {
    return nil
  }
  if err != nil {
    return err
  }

  return json.Unmarshal(data, dest)
}
