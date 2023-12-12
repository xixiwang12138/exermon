package lock

import (
	"context"
	"strconv"
	"time"

	"github.com/xixiwang12138/exermon/redis"
)

const (
	DefaultAutoUnlockTime = 2 * time.Second
	Prefix                = "L:"
)

const (
	//重复上锁的时候，将锁的过期时间续期
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	unlockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen = 16
)

type RedisLock struct {
	key       string
	duration  time.Duration
	expiredAt time.Time
	owner     string
}

func NewRedisLock(key, own string, duration ...time.Duration) *RedisLock {
	var expire time.Duration
	if duration == nil || len(duration) == 0 {
		expire = DefaultAutoUnlockTime
	} else {
		expire = duration[0]
	}
	return &RedisLock{
		key:       key,
		expiredAt: time.Now().Add(expire),
		duration:  expire,
		owner:     own,
	}
}

func (lock *RedisLock) TryLock(ctx context.Context) (bool, error) {
	resp, err := redis.Component.Client.Eval(ctx, lockCommand, []string{Prefix + lock.key}, []string{
		lock.owner, strconv.FormatInt(lock.duration.Milliseconds(), 10)}).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if resp == nil {
		return false, nil
	}
	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}
	return false, nil
}

func (lock *RedisLock) UnLock(ctx context.Context) (bool, error) {
	resp, err := redis.Component.Client.Eval(ctx, unlockCommand, []string{Prefix + lock.key}, []string{lock.owner}).Result()
	if err != nil {
		return false, err
	}
	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}
