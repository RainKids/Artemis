package po

type CloudAccountCredential struct {
	// 账号所在的项目 (openstack)
	ProjectName string `json:"project_name"`

	// 账号所在的域 (openstack)
	// default: Default
	DomainName string `json:"domain_name"`

	// 用户名 (openstack, zstack, esxi)
	Username string `json:"username"`

	// 密码 (openstack, zstack, esxi)
	Password string `json:"password"`

	// 认证地址 (openstack,zstack)
	AuthUrl string `json:"auto_url"`

	// 秘钥id (Aliyun, Aws, huawei, ucloud, ctyun, zstack, s3)
	AccessKeyId string `json:"access_key_id"`

	// 秘钥key (Aliyun, Aws, huawei, ucloud, ctyun, zstack, s3)
	AccessKeySecret string `json:"access_key_secret"`

	// 环境 (Azure, Aws, huawei, ctyun, aliyun)
	Environment string `json:"environment"`

	// 目录ID (Azure)
	DirectoryId string `json:"directory_id"`

	// 客户端ID (Azure)
	ClientId string `json:"client_id"`

	// 客户端秘钥 (Azure)
	ClientSecret string `json:"client_secret"`

	// 主机IP (esxi)
	Host string `json:"host"`

	// 主机端口 (esxi)
	Port int `json:"port"`

	// 端点 (s3) 或 Apsara(飞天)
	Endpoint string `json:"endpoint"`

	// app id (Qcloud)
	AppId string `json:"app_id"`

	//秘钥ID (Qcloud)
	SecretId string `json:"secret_id"`

	//秘钥key (Qcloud)
	SecretKey string `json:"secret_key"`

	// 飞天允许的最高组织id, 默认为1
	OrganizationId int `json:"organization_id"`

	// Google服务账号email (gcp)
	GCPClientEmail string `json:"gcp_client_email"`
	// Google服务账号project id (gcp)
	GCPProjectId string `json:"gcp_project_id"`
	// Google服务账号秘钥id (gcp)
	GCPPrivateKeyId string `json:"gcp_private_key_id"`
	// Google服务账号秘钥 (gcp)
	GCPPrivateKey string `json:"gcp_private_key"`

	// 默认区域Id, Apara及HCSO需要此参数
	// example: cn-north-2
	// required: true
	DefaultRegion string `default:"$DEFAULT_REGION" metavar:"$DEFAULT_REGION"`

	// Huawei Cloud Stack Online
	*HCSOEndpoints
}

type CloudAccount struct {
	// 账号信息，各个平台字段不尽相同，以下是各个平台账号创建所需要的字段
	//
	//
	//
	// | 云平台     |字段                | 翻译              | 是否必传  | 默认值    | 可否更新      | 获取方式   |
	// | ------     |------              | ------            | --------- | --------  |--------       |--------    |
	// |Aliyun      |access_key_id       |秘钥ID             | 是        |            |    是        |            |
	// |Aliyun      |access_key_secret   |秘钥Key            | 是        |            |    是        |            |
	// |Qcloud      |app_id              |APP ID             | 是        |            |    否        |            |
	// |Qcloud      |secret_id           |秘钥ID             | 是        |            |    是        |            |
	// |Qcloud      |secret_key          |秘钥Key            | 是        |            |    是        |            |
	// |OpenStack   |project_name        |用户所在项目       | 是        |            |    是        |            |
	// |OpenStack   |username            |用户名             | 是        |            |    是        |            |
	// |OpenStack   |password            |用户密码           | 是        |            |    是        |            |
	// |OpenStack   |auth_url            |认证地址           | 是        |            |    否        |            |
	// |OpenStack   |domain_name         |用户所在的域       | 否        |Default     |    是        |            |
	// |VMware      |username            |用户名             | 是        |            |    是        |            |
	// |VMware      |password            |密码               | 是        |            |    是        |            |
	// |VMware      |host                |主机IP或域名       | 是        |            |    否        |            |
	// |VMware      |port                |主机端口           | 否        |443         |    否        |            |
	// |Azure       |directory_id        |目录ID             | 是        |            |    否        |            |
	// |Azure       |environment         |区域               | 是        |            |    否        |            |
	// |Azure       |client_id           |客户端ID           | 是        |            |    是        |            |
	// |Azure       |client_secret       |客户端密码         | 是        |            |    是        |            |
	// |Huawei      |access_key_id       |秘钥ID             | 是        |            |    是        |            |
	// |Huawei      |access_key_secret   |秘钥               | 是        |            |    是        |            |
	// |Huawei      |environment         |区域               | 是        |            |    否        |            |
	// |Aws         |access_key_id       |秘钥ID             | 是        |            |    是        |            |
	// |Aws         |access_key_secret   |秘钥               | 是        |            |    是        |            |
	// |Aws         |environment         |区域               | 是        |            |    否        |            |
	// |Ucloud      |access_key_id       |秘钥ID             | 是        |            |    是        |            |
	// |Ucloud      |access_key_secret   |秘钥               | 是        |            |    是        |            |
	// |Google      |project_id          |项目ID             | 是        |            |    否        |            |
	// |Google      |client_email        |客户端email        | 是        |            |    否        |            |
	// |Google      |private_key_id      |秘钥ID             | 是        |            |    是        |            |
	// |Google      |private_key         |秘钥Key            | 是        |            |    是        |            |
	Account string `json:"account"`

	// swagger:ignore
	Secret string

	// 认证地址
	AccessUrl string `json:"access_url"`
}

type HCSOEndpoints struct {
	caches map[string]string

	// 华为私有云Endpoint域名
	// example: hcso.com.cn
	// required:true
	EndpointDomain string `default:"$HUAWEI_ENDPOINT_DOMAIN" metavar:"$HUAWEI_ENDPOINT_DOMAIN"`

	// 默认DNS
	// example: 10.125.0.26,10.125.0.27
	// required: false
	DefaultSubnetDns string `default:"$HUAWEI_DEFAULT_SUBNET_DNS" metavar:"$HUAWEI_DEFAULT_SUBNET_DNS"`

	// 弹性云服务
	Ecs string `default:"$HUAWEI_ECS_ENDPOINT"`
	// 云容器服务
	Cce string `default:"$HUAWEI_CCE_ENDPOINT"`
	// 弹性伸缩服务
	As string `default:"$HUAWEI_AS_ENDPOINT"`
	// 统一身份认证服务
	Iam string `default:"$HUAWEI_IAM_ENDPOINT"`
	// 镜像服务
	Ims string `default:"$HUAWEI_IMS_ENDPOINT"`
	// 云服务器备份服务
	Csbs string `default:"$HUAWEI_CSBS_ENDPOINT"`
	// 云容器实例 CCI
	Cci string `default:"$HUAWEI_CCI_ENDPOINT"`
	// 裸金属服务器
	Bms string `default:"$HUAWEI_BMS_ENDPOINT"`
	// 云硬盘 EVS
	Evs string `default:"$HUAWEI_EVS_ENDPOINT"`
	// 云硬盘备份 VBS
	Vbs string `default:"$HUAWEI_VBS_ENDPOINT"`
	// 对象存储服务 OBS
	Obs string `default:"$HUAWEI_OBS_ENDPOINT"`
	// 虚拟私有云 VPC
	Vpc string `default:"$HUAWEI_VPC_ENDPOINT"`
	// 弹性负载均衡 ELB
	Elb string `default:"$HUAWEI_ELB_ENDPOINT"`
	// 合作伙伴运营能力
	Bss string `default:"$HUAWEI_BSS_ENDPOINT"`
	// Nat网关 NAT
	Nat string `default:"$HUAWEI_NAT_ENDPOINT"`
	// 分布式缓存服务
	Dcs string `default:"$HUAWEI_DCS_ENDPOINT"`
	// 关系型数据库 RDS
	Rds string `default:"$HUAWEI_RDS_ENDPOINT"`
	// 云审计服务
	Cts string `default:"$HUAWEI_CTS_ENDPOINT"`
	// 监控服务 CloudEye
	Ces string `default:"$HUAWEI_CES_ENDPOINT"`
	// 企业项目
	Eps string `default:"$HUAWEI_EPS_ENDPOINT"`
	// 文件系统
	SfsTurbo string `default:"$HUAWEI_SFS_TURBO_ENDPOINT"`
	// Modelarts
	Modelarts string `default:"$HUAWEI_MODELARTS_ENDPOINT"`
}
