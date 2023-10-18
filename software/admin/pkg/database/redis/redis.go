package redis

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// / ClientType 定义redis client 结构体
type RedisDB struct {
	Client redis.UniversalClient
}

// Client  redis连接类型
var RedisClient RedisDB

// Options redis option
type Options struct {
	URL         []string // host:port
	MaxIdle     int      // 最大空闲连接数
	MaxActive   int      // 最大连接数
	IdleTimeout int      // 空闲连接超时时间
	Timeout     int      // 操作超时时间
	Network     string   // tcp or udp
	Password    string
}

// NewOptions for redis
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("redis", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}

	logger.Info("load redis options success", zap.Any("redis options", o))
	return o, err
}

// New redis pool conn
func New(o *Options) (*RedisDB, error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        o.URL,
		Password:     o.Password, // no password set
		DB:           0,          // use default DB
		MaxIdleConns: o.MaxIdle,
		PoolSize:     o.MaxActive,
		DialTimeout:  time.Duration(o.IdleTimeout) * time.Second,
	})
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, err
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return nil, err
	}
	RedisClient = RedisDB{rdb}
	return &RedisClient, nil
}

// DistributedLock 并发锁
func (r *RedisDB) DistributedLock(ctx context.Context, key string, expire int, value time.Time) (bool, error) {
	exists, err := r.Client.Do(ctx, "set", key, value, "nx", "ex", expire).Result()
	if err != nil {
		return false, errors.New("执行 set nx ex 失败")
	}

	// 锁已存在，已被占用
	if exists != nil {
		return false, nil
	}

	return true, nil
}

// ReleaseLock 释放锁
func (r *RedisDB) ReleaseLock(ctx context.Context, key string) (bool, error) {
	v, err := r.Client.Do(ctx, "DEL", key).Bool()
	return v, err
}

// ReleaseLockWithLua 释放锁 使用lua脚本执行
func (r *RedisDB) ReleaseLockWithLua(ctx context.Context, key string, value time.Time) (int, error) {
	// keyCount表示lua脚本中key的个数
	lua := redis.NewScript(ScriptDeleteLock)
	// lua脚本中的参数为key和value
	res, err := lua.Run(ctx, r.Client, []string{key}, value).Int()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *RedisDB) DoSomething(ctx context.Context, key string, expire int, value time.Time) {
	// 获取锁
	canUse, err := r.DistributedLock(ctx, key, expire, value)
	if err != nil {
		panic(err)
	}
	// 占用锁
	if canUse {
		fmt.Println("start do something ...")
		// 释放锁
		_, err := r.ReleaseLock(ctx, key)
		if err != nil {
			panic(err)
		}
	}
}

func (r *RedisDB) DoSomethingWithLua(ctx context.Context, key string, expire int, value time.Time) {
	// 获取锁
	canUse, err := r.DistributedLock(ctx, key, expire, value)
	if err != nil {
		panic(err)
	}
	// 占用锁
	if canUse {
		fmt.Println("start do something ...")
		// 释放锁 lua脚本执行原子性删除
		_, err := r.ReleaseLockWithLua(ctx, key, value)
		if err != nil {
			panic(err)
		}
	}
}

// LuaTokenBucket 通过lua脚本实现令牌桶算法限流
func (r *RedisDB) LuaTokenBucket(ctx context.Context, key string, capacity, rate, now int64) (bool, error) {
	lua := redis.NewScript(ScriptTokenLimit)
	// lua脚本中的参数为key和value
	res, err := lua.Run(ctx, r.Client, []string{key}, capacity, rate, now).Bool()
	if err != nil {
		return false, err
	}
	return res, nil
}

// UnionStore 合并zset的key
func (r *RedisDB) UnionStore(ctx context.Context, rankDays int, keyRank string, c redis.Conn) error {
	today := time.Now()
	unionKeys := make([]string, 0, rankDays+3)
	unionKeys = append(unionKeys, keyRank, strconv.Itoa(rankDays))
	for i := 0; i < rankDays; i++ {
		key := fmt.Sprintf("count:%s", today.AddDate(0, 0, -i).Format("20060102"))
		unionKeys = append(unionKeys, key)
	}

	// 合并一周
	_, err := r.Client.ZUnionStore(ctx, "", &redis.ZStore{Keys: unionKeys}).Result()

	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisDB) Delete(ctx context.Context, key string) (bool, error) {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rdb *RedisDB) HDel(c context.Context, key string, fields ...string) error {

	return rdb.Client.HDel(c, key, fields...).Err()
}

func (rdb *RedisDB) HSan(c context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rdb.Client.HScan(c, key, cursor, match, count).Result()
}

// ProviderSet inject redis settings
var ProviderSet = wire.NewSet(New, NewOptions)
