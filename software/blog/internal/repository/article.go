package repository

import (
	"blog/internal/biz/dto"
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"blog/pkg/database/es"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func queryBuild(keyword string) *elastic.BoolQuery {
	var BoolQuery *elastic.BoolQuery
	titleQuery := elastic.NewMatchQuery("title", keyword)
	contentQuery := elastic.NewMatchQuery("content", keyword)
	return BoolQuery.Should(titleQuery, contentQuery)
}

type articleRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
	es     *es.Client
}

func newArticleRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB, es *es.Client) ArticleRepository {
	return &articleRepository{
		logger: logger.With(zap.String("type", "ArticleRepository")),
		db:     db.Postgres,
		rdb:    rdb,
		es:     es,
	}
}

// IndexExists 查索引是否存在
func (a *articleRepository) IndexExists(article *po.Article) bool {
	context.Background()
	exists, err := a.es.
		IndexExists(context.Background(), article.Index(), true)
	if err != nil {
		a.logger.Error(err.Error())
		return exists
	}
	return exists
}

// CreateIndex 创建索引
func (a *articleRepository) CreateIndex(article *po.Article) error {
	if a.IndexExists(article) {
		// 有索引
		err := a.es.RemoveIndex(context.Background(), article.Index())
		if err != nil {
			return err
		}
		return nil
	}
	// 无索引
	// 创建索引
	err := a.es.CreateIndex(context.Background(), article.Index(), article.Mapping(), true)
	if err != nil {
		a.logger.Error("创建索引失败", zap.Error(err))
		return err
	}
	a.logger.Info(fmt.Sprintf("索引 %s 创建成功", article.Index()))
	return nil

}

func (a *articleRepository) Migrate() error {
	article := new(po.Article)
	err := a.db.AutoMigrate(article)
	err = a.CreateIndex(article)
	return err
}

func (a *articleRepository) Search(search *dto.ArticleSearchParams) (*vo.ArticleSearchList, error) {
	res, err := a.es.Client.Search(po.Article{}.Index()).
		Query(queryBuild(search.Keyword)).
		Highlight(elastic.NewHighlight().Field("title").PreTags("<font color='red'>").PostTags("</font>")).
		Highlight(elastic.NewHighlight().Field("content").PreTags("<font color='green'>").PostTags("</font>")).
		DefaultRescoreWindowSize(10000).
		Sort("id", false).
		From(search.Page).
		Size(search.Page).
		Do(context.Background())
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}
	count := res.Hits.TotalHits.Value
	list := make([]*vo.ArticleSearch, 0, count)
	for _, hit := range res.Hits.Hits {
		var article vo.ArticleSearch
		err = json.Unmarshal(hit.Source, &article)
		if err != nil {
			a.logger.Error(err.Error())
			continue
		}
		list = append(list, &vo.ArticleSearch{ID: article.ID,
			Title:   hit.Highlight["title"][0],
			Content: hit.Highlight["content"][0]})
	}
	return &vo.ArticleSearchList{Result: list, Count: count}, nil
}
