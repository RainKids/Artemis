package po

import (
	"blog/global"
	"blog/pkg/database"
)

type Article struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`   // 主键ID
	Title    string `json:"title" structs:"title"`                // 文章标题
	Keyword  string `json:"keyword,omit(list)" structs:"keyword"` // 关键字
	Abstract string `json:"abstract" structs:"abstract"`          // 文章简介
	Content  string `json:"content,omit(list)" structs:"content"` // 文章内容

	LookCount    int    `json:"lookCount" structs:"lookCount"`       // 浏览量
	CommentCount int    `json:"commentCount" structs:"commentCount"` // 评论量
	LikeCount    int    `json:"LikeCount" structs:"LikeCount"`       // 点赞量
	CollectCount int    `json:"collectCount" structs:"collectCount"` // 收藏量
	UserID       uint   `json:"userId" structs:"userId"`             // 用户id
	Category     string `json:"category" structs:"category"`         // 文章分类
	Source       string `json:"source" structs:"source"`             // 文章来源
	Link         string `json:"link" structs:"link"`                 // 原文链接
	BannerID     uint   `json:"bannerId" structs:"bannerId"`         // 文章封面ID
	BannerUrl    string `json:"bannerUrl" structs:"bannerUrl"`       // 文章封面

	Tags database.Array `json:"tags" structs:"tags"` // 文章标签
	global.OperateBy
	global.ModelTime
}

func (Article) Mapping() string {
	return `
{
  "settings": {
    "analysis": {
      "analyzer": {
        "text_analyzer": {
          "tokenizer": "ik_max_word"
        }
      }
    }
  },
  "mappings": {
    "properties": {
 	"all": {
        "type": "text",
        "analyzer": "text_analyzer"
      },
      "title": { 
        "type": "text",
 		"analyzer": "text_analyzer",
		"copy_to": "all"
      },
      "keyword": { 
        "type": "keyword"
      },
      "id": {
        "type": "integer"
      },
      "abstract": { 
        "type": "text"
      },
      "content": { 
        "type": "text",
		"analyzer": "text_analyzer",
		"copy_to": "all"
      },
      "look_count": {
        "type": "integer"
      },
      "comment_count": {
        "type": "integer"
      },
      "like_count": {
        "type": "integer"
      },
      "collects_count": {
        "type": "integer"
      },
      "user_id": {
        "type": "integer"
      },
      "category_id": { 
        "type": "integer"
      },
      "source": { 
        "type": "keyword"
      },
      "link": { 
        "type": "keyword"
      },
      "banner_id": {
        "type": "integer"
      },
      "banner_url": { 
        "type": "keyword"
      },
      "tags": { 
        "type": "keyword"
      },
      "created_at":{
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      },
      "updated_at":{
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      }
    }
  }
}
`
}

func (Article) Index() string {
	return "article_index"
}

func (Article) TableName() string {
	return "blog_article"
}
