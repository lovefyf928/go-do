package redis

import (
	"context"
	"github.com/go-redis/redis"
	"go-do/common/conf"
)

var ctx = context.Background()

var Rdb *redis.Client

func init() {
	newRdb := redis.NewClient(&redis.Options{
		Addr:     conf.ConfigInfo.DataSource.Redis.Url,
		Password: conf.ConfigInfo.DataSource.Redis.Passwd, // no password set
		DB:       conf.ConfigInfo.DataSource.Redis.Db,     // use default DB
	})

	Rdb = newRdb

}
