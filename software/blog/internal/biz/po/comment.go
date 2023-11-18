package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"` // 主键ID
	ObjId     int64              `json:"objId"`         //评论对象
	Reply     int64              `json:"replyID"`       //回复id
	Type      int8               `json:"type"`          //对象类型
	User      int64              `json:"userId"`        //作者用户id
	Root      primitive.ObjectID `json:"root"`          //根评论id不为nil是回复评论
	Parent    primitive.ObjectID `json:"parent"`        //父评论id不为nil是root评论
	Floor     int64              `json:"floor"`         //评论楼层
	Count     int64              `json:"count"`         //评论总数
	RootCount int64              `json:"rootCount"`     //根评论总数
	Like      int64              `json:"like"`          //点赞数
	Hate      int64              `json:"hate"`          //点踩数
	Status    int8               `json:"status"`        //状态0-正常 1-隐藏
	Attrs     int32              `json:"attrs"`         //属性0-运营置顶1-对象用户置顶2-大数据过滤
	Platform  int8               `json:"platform"`      //平台类型
	Device    string             `json:"device"`        // 设备信息
	Ip        string             `json:"ip"`            // 评论Ip
	Message   string             `json:"content"`       // 评论内容
	Meta      string             `json:"meta"`          // 评论元数据
	CreatedAt time.Time          `json:"createdAt" `    // 创建时间
	UpdatedAt time.Time          `json:"updatedAt"`     // 更新时间
}
type CommentUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"` // 主键ID
	ObjID    int64              `json:"objId"`         //对象id
	Username string             `json:"username"`
}

func (Comment) TableName() string {
	return "blog_comment"
}

func (CommentUser) TableName() string {
	return "blog_comment_user"
}
