package service

import (
	"admin/global"
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/internal/repository"
	"context"
	"go.uber.org/zap"
)

type postService struct {
	logger         *zap.Logger
	postRepository repository.PostRepository
}

func newPostService(logger *zap.Logger, postRepository repository.PostRepository) PostService {
	return &postService{
		logger:         logger.With(zap.String("type", "PostService")),
		postRepository: postRepository,
	}
}

func (p *postService) List(c context.Context, params *dto.PostSearchParams) (*vo.PostList, error) {
	list, count, err := p.postRepository.List(c, params)
	if err != nil {
		return nil, err
	}
	return &vo.PostList{
		Result: list,
		Count:  count,
	}, nil
}

func (p *postService) Create(c context.Context, req *dto.PostRequest) (*po.Post, error) {
	return p.postRepository.Create(c, &po.Post{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
		OperateBy: global.OperateBy{CreateBy: req.OperateBy},
	})
}
func (p *postService) Retrieve(c context.Context, id int64) (*po.Post, error) {
	return p.postRepository.Retrieve(c, id)
}
func (p *postService) Update(c context.Context, id int64, req *dto.PostRequest) error {
	return p.postRepository.Update(c, &po.Post{
		ID:        id,
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
		OperateBy: global.OperateBy{UpdateBy: req.OperateBy},
	})
}
func (p *postService) Delete(c context.Context, id int64) error {
	return p.postRepository.Delete(c, id)
}
