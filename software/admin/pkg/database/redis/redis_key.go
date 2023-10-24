package redis

const (
	// KeyAdminLoginFailedCount 1小时用户登陆失败数key blog:login:failed:count:2
	KeyAdminLoginFailedCount = "blog:login:failed:count:%s"
	// KeyArticleCount 24小时文章阅读数key blog:article:count:20201020
	KeyArticleCount = "blog:article:read:count:%s"
)
