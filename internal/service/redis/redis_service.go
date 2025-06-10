package redis

import (
	"chat_server/internal/config"
	"chat_server/pkg/zlog"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.Db
	addr := host + ":" + strconv.Itoa(port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// 设置有超时时间的 key
func SetKeyEx(key string, value string, timeout time.Duration) error {
	err := redisClient.Set(ctx, key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取 key，如果不存在，err = nil
func GetKey(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info("该 key 不存在")
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// 获取 key，如果不存在，err != nil
func GetKeyNilIsErr(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 获取有 prefix 前缀的 key，如果不存在，err != nil
func GetKeyWithPrefixNilIsErr(prefix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return "", err
		}

		if len(keys) == 0 {
			zlog.Info("没有找到相关前缀 key")
			return "", redis.Nil
		}

		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关前缀 key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于 1 的 key, 查找异常")
			return "", errors.New("找到了数量大于 1 的 key, 查找异常")
		}
	}
}

// 获取有 suffix 后缀的 key，如果不存在，err != nil
func GetKeyWithSuffixNilIsErr(suffix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
		if err != nil {
			return "", err
		}

		if len(keys) == 0 {
			zlog.Info("没有找到相关后缀 key")
			return "", redis.Nil
		}

		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关后缀 key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于 1 的 key, 查找异常")
			return "", errors.New("找到了数量大于 1 的 key, 查找异常")
		}
	}
}

// 如果 key 存在，删除该 key
func DelKeyIfExists(key string) error {
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exists == 1 { // 键存在
		delErr := redisClient.Del(ctx, key).Err()
		if delErr != nil {
			return delErr
		}
	}
	// 无论键是否存在，都不返回错误
	return nil
}

// 如果有该 key，则删除
func DelKeyWithPattern(pattern string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			zlog.Info("没有找到对应 key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			zlog.Info(fmt.Sprintln("成功删除相关对应的 key", keys))
		}
	}

	return nil
}

// 如果有该前缀的 key，则删除
func DelKeyWithPrefix(prefix string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			zlog.Info("没有找到相关前缀 key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			zlog.Info(fmt.Sprintln("成功删除相关前缀 key", keys))
		}
	}

	return nil
}

// 如果有该后缀的 key，则删除
func DelKeyWithSuffix(suffix string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			zlog.Info("没有找到相关后缀 key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			zlog.Info(fmt.Sprintln("成功删除相关后缀 key", keys))
		}
	}

	return nil
}

// 删除所有的 key
func DeleteAllRedisKeys() error {
	var cursor uint64 = 0

	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		cursor = nextCursor

		if len(keys) > 0 {
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
