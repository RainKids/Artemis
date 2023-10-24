package captcha

import (
	"admin/pkg/database/redis"
	"admin/pkg/logger"
	"context"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"time"
)

type RedisStore struct {
	Expiration int
	PreKey     string
	Context    context.Context
}

func NewDefaultRedisStore() *RedisStore {
	return &RedisStore{
		Expiration: CaptchaExpireTime,
		PreKey:     "CAPTHCA_",
	}
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.Context = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	c := context.Background()
	err := redis.RedisClient.Set(c, rs.PreKey+id, value, time.Duration(int64(rs.Expiration))*time.Second)
	if err != nil {
		logger.Logger.Error("RedisStoreSetError", zap.Error(err))
		return err
	}
	return err
}

func (rs *RedisStore) Get(key string, clear bool) string {
	c := context.Background()
	val, err := redis.RedisClient.Get(c, key)
	if err != nil {
		logger.Logger.Error("RedisStoreGetError", zap.Error(err))
		return ""
	}
	if clear {
		res, err := redis.RedisClient.Delete(c, key)
		if !res && err != nil {
			logger.Logger.Error("RedisStoreGetError", zap.Error(err))
			return ""
		}
	}
	return val
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	key := rs.PreKey + id
	v := rs.Get(key, clear)
	return v == answer
}
