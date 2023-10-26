package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/pkg/database"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newPostRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) PostRepository {
	return &postRepository{
		logger: logger.With(zap.String("type", "PostRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (p *postRepository) Migrate() error {
	return p.db.AutoMigrate(&po.Post{})
}

func (p *postRepository) List(c context.Context, params *dto.PostSearchParams) (list []*po.Post,
	count int64, err error) {
	db := p.db.WithContext(c).Model(&po.Post{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Code != "" {
		db = db.Where("code LIKE ?", "%"+params.Code+"%")
	}
	db = db.Where("status = ?", params.Status)
	if err = db.Scopes(database.Paginate(params.Page, params.PageSize)).Find(&list).Count(&count).Error; err != nil {
		return nil, 0, errors.Errorf("get post list db error: %s", err)
	}
	return
}

func (p *postRepository) Retrieve(c context.Context, id int64) (*po.Post, error) {
	var post po.Post
	err := p.db.WithContext(c).Where("id = ?", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%s]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	return &post, nil
}

func (p *postRepository) Create(c context.Context, post *po.Post) (*po.Post, error) {
	if err := p.db.WithContext(c).Create(post).Error; err != nil {
		return nil, errors.Errorf("create post db error: %s", err)
	}
	return post, nil

}

func (p *postRepository) Update(c context.Context, post *po.Post) error {
	if err := p.db.WithContext(c).Updates(post).Error; err != nil {
		return errors.Errorf("update role db error: %s", err)
	}
	return nil
}

func (p *postRepository) Delete(c context.Context, id int64) error {
	var post po.Post
	db := p.db.WithContext(c).Model(&post).Delete(&post, id)
	err := db.Error
	if err != nil {
		return errors.Errorf("delete post db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete post")
	}
	return nil
}
