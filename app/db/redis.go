package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func RedisConn(ctx context.Context) (*redis.Client, error) {
    conn := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
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
