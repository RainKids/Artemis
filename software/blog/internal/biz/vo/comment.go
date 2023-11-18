package vo

import (
	"blog/internal/biz/po"
	"time"
)

type CommentUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type ArticleCommentList struct {
	Result []*ArticleComment `json:"result"`
	Count  int64             `json:"count"`
}

type ArticleComment struct {
	ID        string                `bson:"_id" json:"id"` // 主键ID
	Users     []po.CommentUser      `json:"users,omitempty"`
	Replies   []po.CommentUser      `json:"replies,omitempty"`
	User      CommentUser           `json:"user"`
	Reply     CommentUser           `json:"reply"`
	Ip        string                `json:"ip"`         // 评论Ip
	Message   string                `json:"content"`    // 评论内容
	CreatedAt time.Time             `json:"createdAt" ` // 创建时间
	Children  []ArticleReplyComment `json:"children"`
}

type ArticleReplyComment struct {
	ID        string           `bson:"_id" json:"id"` // 主键ID
	Users     []po.CommentUser `json:"users,omitempty"`
	Replies   []po.CommentUser `json:"replies,omitempty"`
	User      CommentUser      `json:"user"`
	Reply     CommentUser      `json:"reply"`
	Ip        string           `json:"ip"`         // 评论Ip
	Message   string           `json:"content"`    // 评论内容
	CreatedAt time.Time        `json:"createdAt" ` // 创建时间
}

type ArticleReplyCommentList struct {
	Result []*ArticleComment `json:"result"`
	Count  int64             `json:"count"`
}
