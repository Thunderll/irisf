package redis_manager

import (
	"context"
	"iris_project_foundation/config"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	var (
		err error
	)
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.GConfig.Redis.Addr,
		Password: config.GConfig.Redis.Password,
		DB:       int(config.GConfig.Redis.DB),
	})

	if _, err = RDB.Ping(context.TODO()).Result(); err != nil {
		log.Fatalf("[ERR] redis连接失败.\n%v", err)
	}
	log.Println("[INFO] redis连接成功.")
}

func LockKey(key string) bool {
	return RDB.SetNX(context.TODO(), key, 1, 10*time.Second).Val()
}

func Unlock(key string) {
	RDB.Del(context.TODO(), key)
}
