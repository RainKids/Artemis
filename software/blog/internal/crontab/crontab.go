package crontab

import (
	"blog/internal/biz/po"
	"blog/pkg/database/es"
	"blog/pkg/database/mongo"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gochore/dcron"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"strings"
	"time"
)

// DefaultCronJobService CronJob模块service默认对象
type DefaultCronJobService struct {
	logger *zap.Logger
	v      *viper.Viper
	rdb    *redis.RedisDB
	es     *es.Client
	db     *postgres.DB
	mongo  *mongo.MongoDB
}

// NewDefaultCronJobService 初始化
func NewDefaultCronJobService(
	logger *zap.Logger,
	v *viper.Viper,
	rdb *redis.RedisDB,
	es *es.Client,
	db *postgres.DB,
	mongo *mongo.MongoDB,
) *DefaultCronJobService {
	return &DefaultCronJobService{
		logger: logger.With(zap.String("type", "DefaultCronJobService")),
		v:      v,
		rdb:    rdb,
		es:     es,
		db:     db,
		mongo:  mongo,
	}
}

// RedisToES 同步redis数据到ES中 数据持久化
func (s *DefaultCronJobService) RedisToES(c context.Context) error {
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(redis.KeyArticleCount, todayStr)
	// zrange Key 0 -1 从redis取出zset每天所有的文章id及浏览次数
	articles, err := s.rdb.Client.ZRangeWithScores(context.Background(), key, 0, -1).Result()
	if err != nil {
		s.logger.Error("ZRANGE", zap.Any("error", err))
		return err
	}
	for _, val := range articles {
		IdStr := val.Member.(string)
		countStr := convertor.ToString(val.Score)
		if err != nil {
			s.logger.Warn("Article:strconv.Atoi count failed", zap.Any("error", err))
			continue
		}

		// 更新es的文章浏览次数
		Doc := make(map[string]string)
		Doc["count"] = countStr
		_ = s.es.Update(context.Background(), "article", IdStr, "", Doc)
	}
	return nil
}

// RedisToMongo 同步redis数据到Mongo中 数据持久化
func (s *DefaultCronJobService) RedisToMongo(c context.Context) error {
	//TODO implement me 文章评论数信息同步
	if _, ok := dcron.TaskFromContext(c); ok {
		keys, err := s.rdb.Client.Keys(c, "blog:article:comment:*").Result()
		if err != nil {
			return err
		}
		for i := 0; i < len(keys); i++ {
			keyArray := strings.Split(keys[i], ":")
			commentID := keyArray[len(keyArray)-1]
			col := s.mongo.DB.Collection(po.Comment{}.TableName())
			var comment po.Comment
			err := col.FindOne(context.TODO(), bson.D{{"_id", commentID}}).Decode(&comment)
			if err != nil {
				return err
			}
			likeCount, err := s.rdb.HGet(c, keys[i], "like").Int64()
			if err != nil {
				s.rdb.HSet(c, keys[i], "like", comment.Like)
			}
			hateCount, err := s.rdb.HGet(c, keys[i], "hate").Int64()
			if err != nil {
				s.rdb.HSet(c, keys[i], "hate", comment.Hate)
			}
			floorCount, err := s.rdb.HGet(c, keys[i], "floor").Int64()
			if err != nil {
				s.rdb.HSet(c, keys[i], "floor", comment.Floor)
			}
			commentCount, err := s.rdb.HGet(c, keys[i], "count").Int64()
			if err != nil {
				s.rdb.HSet(c, keys[i], "comment", comment.Count)
			}
			_, err = col.UpdateOne(c, bson.M{"_id": commentID}, bson.D{{"like", likeCount},
				{"hate", hateCount}, {"floor", floorCount}, {"count", commentCount}})
			s.logger.Error("update comment error", zap.Error(err))
		}
	}
	return nil
}

// RedisToDB 同步redis数据到DB中 数据持久化
func (s *DefaultCronJobService) RedisToPostgres(c context.Context) error {
	//TODO implement me  文章点赞数
	return nil

}

// PostgresToEs 同步Postgres数据到Es中 数据检索
func (s *DefaultCronJobService) PostgresToEs(c context.Context) error {
	//TODO implement me  文章搜索
	return nil

}
