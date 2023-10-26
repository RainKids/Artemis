package po

import (
	"admin/global"
	"time"
)

type User struct {
	ID            int64     `json:"id" form:"id" gorm:"primarykey;AUTO_INCREMENT"`
	Username      string    `json:"username" form:"username" gorm:"column:user_name;comment:用户昵称;type:varchar(50);"`
	Phone         string    `json:"phone" form:"Phone" gorm:"uniqueIndex;column:phone;comment:手机号;type:varchar(11);"`
	Img           string    `json:"img" form:"img" gorm:"column:img;type:varchar(255);"`
	Password      string    `json:"password" form:"password" gorm:"column:password;comment:MD5加密后的密码;type:varchar(128);"`
	IntroduceSign string    `json:"introduceSign" form:"introduceSign" gorm:"column:introduce_sign;comment:个性签名;type:varchar(100);"`
	Status        int       `json:"status" form:"status" gorm:"column:status;comment:锁定标识字段(0-未锁定 1-已锁定);type:int2"`
	LockedFlag    int       `json:"lockedFlag" form:"lockedFlag" gorm:"column:locked_flag;comment:锁定标识字段(0-未锁定 1-已锁定);type:int2"`
	Email         string    `json:"email" gorm:"column:email;comment:邮箱;type:varchar(128)"` //注册邮箱
	Sex           string    `json:"sex" gorm:"column:sex;comment:性别m男g女;type:varchar(2)"`   //性别，m男，g女                                                   //
	Birthday      string    `json:"birthday"`                                               //出生年月日
	QQ            string    `json:"qq"`                                                     //QQ号码
	Wechat        string    `json:"wechat"`                                                 //微信
	UserType      int       `json:"type" gorm:"default:0"`                                  //用户类型
	Salt          string    `json:"salt" gorm:"column:salt;comment:密码加盐;type:varchar(32)"`  //密码，加密存储
	Description   string    `json:"description" gorm:"size:255;comment:描述"`                 // 描述
	RoleId        int64     `json:"roleId" gorm:"size:20;comment:角色ID"`
	DeptId        int64     `json:"deptId" gorm:"size:20;comment:部门"`
	PostId        int64     `json:"postId" gorm:"size:20;comment:岗位"`
	LastLoginTime time.Time `json:"lastLoginTime"` //最后登录时间
	LastLoginIP   string    `json:"lastLoginIP"`   //最后登录时间
	DeptIds       []int     `json:"deptIds" gorm:"-"`
	PostIds       []int     `json:"postIds" gorm:"-"`
	RoleIds       []int     `json:"roleIds" gorm:"-"`
	Dept          *Dept     `json:"dept"`
	global.ModelTime
	global.OperateBy
}

func (User) TableName() string {
	return "sys_user"
}

type UserToken struct {
	Id         int       `json:"id" form:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UserId     int       `json:"userId" form:"userId" gorm:"uniqueIndex"`
	Token      string    `json:"token" form:"token" gorm:"column:token;comment:token值(32位字符串);type:varchar(32);"`
	UpdateTime time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;type:timestamptz"`
	ExpireTime time.Time `json:"expireTime" form:"expireTime" gorm:"column:expire_time;comment:token过期时间;type:timestamptz"`
}

// TableName MallUserToken 表名
func (UserToken) TableName() string {
	return "sys_user_token"
}
