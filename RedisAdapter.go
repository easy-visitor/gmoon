package gmoon

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisAdapter struct {
	*redis.Client
}

func (r RedisAdapter) Name() string {
	return "Redis"
}

func NewRedisAdapter() *RedisAdapter {

	redisCfg := config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		Logger.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		Logger.Info("redis connect ping response:", zap.String("pong", pong))
	}
	return &RedisAdapter{
		Client: client,
	}

}

func Test() {

	s := NewRedisAdapter().Get(context.Background(), "name").String()
	fmt.Println(s)
}
