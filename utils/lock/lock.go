package lock

import (
	"TikTok/config"
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redisClient "github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Lock struct {
	redSync *redsync.Redsync
}

var (
	lockExpiry = time.Second * 10
	retryDelay = time.Millisecond * 100
	tries      = 3
	option     = []redsync.Option{
		redsync.WithExpiry(lockExpiry),
		redsync.WithRetryDelay(retryDelay),
		redsync.WithTries(tries),
	}
	MutexLock Lock
)

func InitLock() error {
	conf := config.Conf.RedSync
	lockClient := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		Password: conf.Pwd,
		PoolSize: conf.PoolSize,
	})
	option[0] = redsync.WithExpiry(time.Second * time.Duration(conf.LockExpire))
	if lockClient == nil {
		panic("failed to call redis.NewClient")
	}
	if _, err := lockClient.Ping(context.Background()).Result(); err != nil {
		panic("Failed to ping redisMutex, err:%s")
	}
	MutexLock.redSync = redsync.New([]redis.Pool{goredis.NewPool(lockClient)}...)
	return nil
}

func (lock *Lock) GetLock(id int64) *redsync.Mutex {
	return lock.redSync.NewMutex(strconv.FormatInt(id, 10), option...)
}

func (lock *Lock) Unlock(mutex *redsync.Mutex) {
	mutex.Unlock()
}
