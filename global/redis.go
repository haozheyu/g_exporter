package global

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb *redis.Client
)

// 初始化连接
func RedisClient(host,port,password string) *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s",host,port),
		Password: password,  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}