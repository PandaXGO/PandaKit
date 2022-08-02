package starter

import (
	"fmt"
	"github.com/XM-GO/PandaKit/config"
	"github.com/XM-GO/PandaKit/logger"
	"github.com/go-redis/redis"
)

func ConnRedis() *redis.Client {
	// 设置redis客户端
	redisConf := config.Conf.Redis
	if redisConf == nil {
		logger.Log.Panic("未找到redis配置信息")
	}
	logger.Log.Infof("连接redis [%s:%d]", redisConf.Host, redisConf.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password, // no password set
		DB:       redisConf.Db,       // use default DB
	})
	// 测试连接
	_, e := rdb.Ping().Result()
	if e != nil {
		logger.Log.Panic(fmt.Sprintf("连接redis失败! [%s:%d]", redisConf.Host, redisConf.Port))
	}
	return rdb
}
