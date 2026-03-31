package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goAi/config"
	"goAi/pkg/logger"
)

var RDB *redis.Client

func InitRedis(cfg config.RedisConfig) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// 测试连接
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("Redis 连接失败: " + err.Error())
	}

	logger.Info("Redis 连接成功")
}
