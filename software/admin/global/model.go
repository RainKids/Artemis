package global

import (
	"gorm.io/gorm"
	"time"
)

type ModelTime struct {
	CreatedAt time.Time `json:"createdAt" ` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`  // 更新时间
}

type DeleteTime struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type OperateBy struct {
	CreateBy int `json:"createBy"`
	UpdateBy int `json:"updateBy"`
}

type Migration struct {
	Version   string    `gorm:"primaryKey"`
	ApplyTime time.Time `gorm:"autoCreateTime"`
}

func (Migration) TableName() string {
	return "sys_migration"
}
