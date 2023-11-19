package repository

import (
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"blog/pkg/database/mongo"
	"blog/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type commentRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
	mongo  *mongo.MongoDB
}

func newCommentRepository(logger *zap.Logger, rdb *redis.RedisDB, mongo *mongo.MongoDB) CommentRepository {
	return &commentRepository{
		logger: logger.With(zap.String("type", "CommentRepository")),
		rdb:    rdb,
		mongo:  mongo,
	}
}

func (t *commentRepository) Migrate() error {
	err := t.mongo.DB.CreateCollection(context.Background(), po.Comment{}.TableName())
	err = t.mongo.DB.CreateCollection(context.Background(), po.CommentUser{}.TableName())
	return err
}

func (t *commentRepository) getUserByID(id int64) (*po.CommentUser, error) {
	col := t.mongo.DB.Collection(po.Comment{}.TableName())
	var user po.CommentUser
	err := col.FindOne(context.TODO(), bson.D{{"objId", id}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (t *commentRepository) Create(c context.Context, comment *po.Comment) (*vo.ArticleReplyComment, error) {

	user, err := t.getUserByID(comment.User)
	if err != nil {
		return nil, err
	}
	reply, err := t.getUserByID(comment.Reply)
	if err != nil {
		return nil, err
	}
	insertOneResult, err := t.mongo.DB.Collection(po.Comment{}.TableName()).InsertOne(c, comment)
	if err != nil {
		return nil, err
	}
	return &vo.ArticleReplyComment{
		ID: insertOneResult.InsertedID.(string),
		User: &vo.CommentUser{
			ID:       user.ObjID,
			Username: user.Username,
		},
		Reply: &vo.CommentUser{
			ID:       reply.ObjID,
			Username: reply.Username,
		},
		Ip:      comment.Ip,
		Message: comment.Message,
	}, nil
}

func (t *commentRepository) ListByArticleID(c context.Context, articleId, offset, limit int64) ([]*vo.ArticleComment, int64, error) {
	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: []bson.M{
					{"parent": primitive.NilObjectID},
					{"root": primitive.NilObjectID},
				}},
				{Key: "objId", Value: articleId},
				{Key: "status", Value: 0},
			}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "blog_comment_user"},
			{Key: "localField", Value: "user"},
			{Key: "foreignField", Value: "objId"},
			{Key: "as", Value: "users"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "blog_comment_user"},
			{Key: "localField", Value: "reply"},
			{Key: "foreignField", Value: "objId"},
			{Key: "as", Value: "replies"},
		}}},
	}
	col := t.mongo.DB.Collection(po.Comment{}.TableName())
	cursor, err := col.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, 0, errors.Errorf("reply comment db error: %s", err)
	}
	count, err := col.CountDocuments(c, pipeline)
	if err != nil {
		return nil, 0, errors.Errorf("get count err: %s", err)
	}
	var comments []*vo.ArticleComment
	if err := cursor.All(context.TODO(), &comments); err != nil {
		return nil, 0, errors.Errorf("root comment to struct err: %s", err)
	}
	//查询回复评论
	for _, comment := range comments {
		var children []vo.ArticleReplyComment
		pipeline = bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "parent", Value: comment.ID}}}},
			bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
			bson.D{{Key: "$skip", Value: 0}},
			bson.D{{Key: "$limit", Value: 3}},
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "blog_comment_user"},
				{Key: "localField", Value: "user"},
				{Key: "foreignField", Value: "objId"},
				{Key: "as", Value: "users"},
			}}},
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "blog_comment_user"},
				{Key: "localField", Value: "reply"},
				{Key: "foreignField", Value: "objId"},
				{Key: "as", Value: "replies"},
			}}},
		}
		cursor, err = col.Aggregate(context.TODO(), pipeline)
		if err != nil {
			return nil, count, errors.Errorf("db reply comment error: %s", err)
		}
		err = cursor.All(context.TODO(), &children)
		if err != nil {
			return nil, count, errors.Errorf("reply comment to struct err: %s", err)
		}
		if len(comment.Users) > 0 {
			comment.User = &vo.CommentUser{
				ID:       comment.Users[0].ObjID,
				Username: comment.Users[0].Username,
			}
		}
		if len(comment.Replies) > 0 {
			comment.Reply = &vo.CommentUser{
				ID:       comment.Replies[0].ObjID,
				Username: comment.Replies[0].Username,
			}
		}
		for _, child := range comments {
			if len(child.Users) > 0 {
				child.User = &vo.CommentUser{
					ID:       child.Users[0].ObjID,
					Username: child.Users[0].Username,
				}
			}
			if len(child.Replies) > 0 {
				child.Reply = &vo.CommentUser{
					ID:       child.Replies[0].ObjID,
					Username: child.Replies[0].Username,
				}
			}
		}
		comment.Children = children
	}
	return comments, count, nil
}

func (t *commentRepository) ListByParentID(c context.Context, ParentID primitive.ObjectID, offset, limit int64) ([]*vo.ArticleReplyComment, int64, error) {
	var reply []*vo.ArticleReplyComment
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "parent", Value: ParentID}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "blog_comment_user"},
			{Key: "localField", Value: "user"},
			{Key: "foreignField", Value: "objId"},
			{Key: "as", Value: "users"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "blog_comment_user"},
			{Key: "localField", Value: "reply"},
			{Key: "foreignField", Value: "objId"},
			{Key: "as", Value: "replies"},
		}}},
	}
	col := t.mongo.DB.Collection(po.Comment{}.TableName())
	cursor, err := col.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, 0, errors.Errorf("reply comment db error: %s", err)
	}
	count, err := col.CountDocuments(c, pipeline)
	if err != nil {
		return nil, 0, errors.Errorf("get count err: %s", err)
	}
	err = cursor.All(context.TODO(), &reply)
	if err != nil {
		return nil, count, errors.Errorf("reply comment to struct err: %s", err)
	}
	for _, comment := range reply {
		if len(comment.Users) > 0 {
			comment.User = &vo.CommentUser{
				ID:       comment.Users[0].ObjID,
				Username: comment.Users[0].Username,
			}
		}
		if len(comment.Replies) > 0 {
			comment.Reply = &vo.CommentUser{
				ID:       comment.Replies[0].ObjID,
				Username: comment.Replies[0].Username,
			}
		}
	}
	return reply, count, nil
}

func (t *commentRepository) Retrieve(c context.Context, id primitive.ObjectID) (*vo.ArticleReplyComment, error) {
	col := t.mongo.DB.Collection(po.Comment{}.TableName())
	var comment po.Comment
	err := col.FindOne(c, bson.D{{"_id", id}}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	user, err := t.getUserByID(comment.User)
	if err != nil {
		return nil, err
	}
	reply, err := t.getUserByID(comment.Reply)
	if err != nil {
		return nil, err
	}
	return &vo.ArticleReplyComment{
		ID: comment.ID.Hex(),
		User: &vo.CommentUser{
			ID:       user.ObjID,
			Username: user.Username,
		},
		Reply: &vo.CommentUser{
			ID:       reply.ObjID,
			Username: reply.Username,
		},
		Ip:        comment.Ip,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (t *commentRepository) Delete(c context.Context, id primitive.ObjectID) error {
	_, err := t.mongo.DB.Collection(po.Comment{}.TableName()).DeleteOne(c, bson.D{{"objId", id}})
	return err
}
