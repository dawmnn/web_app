package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

//初始化连接
//
//func Init() (err error) {
//	client = redis.NewClient(&redis.Options{
//		Addr: fmt.Sprintf("%s:%d",
//			viper.GetString("redis.host"),
//			viper.GetInt("redis.port"),
//		),
//		Password: viper.GetString("redis.password"),
//		DB:       viper.GetInt("redis.db"),
//		PoolSize: viper.GetInt("redis.pool_size"),
//	})
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	_, err = client.Ping(ctx).Result()
//
//	return err
//}

func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()

	return err
}

func Close() {
	_ = client.Close()
}

//func Init(ctx context.Context) {
//	// 初始化 Redis
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: "127.0.0.1:6379",
//	})
//
//	if _, err := redisClient.Ping(ctx).Result(); err != nil {
//		log.Fatalf("init redis failed, err: %v", err)
//	}
//	log.Println("Redis connected successfully")
//
//}
