package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var ctx = context.TODO()

var rdb = redis.NewClient(&redis.Options{
	Addr: "instabug_redis:6379",
})

var mycache = cache.New(&cache.Options{
	Redis:      rdb,
	LocalCache: cache.NewTinyLFU(1000, time.Hour),
})

var Ring = redis.NewRing(&redis.RingOptions{
	Addrs: map[string]string{
		"server1": "instabug_redis:6379",
	},
})

type RedisObj struct {
	Num int64
}

func SaveInRedis(key string, value int64) {
	// set key to value in redis v8
	obj := &RedisObj{
		Num: value,
	}

	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}
}

func GetFromRedis(key string) int64 {
	var res RedisObj
	if err := mycache.Get(ctx, key, &res); err == nil {
		return res.Num
	}
	return -1
}
