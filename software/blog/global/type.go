package global

type ImageType int

const (
	Local   ImageType = 1 // 本地
	QiNiu   ImageType = 2 // 七牛云
	AliYun  ImageType = 3
	Tencent ImageType = 4
	Huawei  ImageType = 5
)

func (s ImageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s ImageType) String() string {
	switch s {
	case Local:
		return "本地"
	case QiNiu:
		return "七牛云"
	case AliYun:
		return "阿里云"
	case Tencent:
		return "腾讯云"
	case Huawei:
		return "华为云"
	default:
		return "其他"
	}
}

type SignStatus int

const (
	SignEmail  SignStatus = 1 //邮箱
	SignPhone  SignStatus = 2
	SignQQ     SignStatus = 3 // QQ
	SignWechat SignStatus = 4
	SignGitee  SignStatus = 5 // Gitee
	SigGithub  SignStatus = 6
)

func (s SignStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s SignStatus) String() string {
	switch s {
	case SignQQ:
		return "QQ"
	case SignGitee:
		return "Gitee"
	case SignEmail:
		return "邮箱"
	case SignWechat:
		return "微信"
	case SigGithub:
		return "Github"
	case SignPhone:
		return "手机"
	default:
		return "其他"
	}
}
