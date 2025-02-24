package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func RedisConn(ctx context.Context) (*redis.Client, error) {
    redisHost := os.Getenv("APP_REDIS_HOST")
    conn := redis.NewClient(&redis.Options{
        Addr: redisHost,
        Password: "",
        DB: 0,
    })
    pong, err := conn.Ping(ctx).Result()
    if err != nil {
        fmt.Println(pong, err, "Connection result")
        return conn, err
    }

    return conn, nil
}
