package conf

type SystemConfig struct {
	Port                int    `mapstructure:"port"`                  // 端口
	BaseURL             string `mapstructure:"base_url"`              // 服务base_url
	DefaultLang         string `mapstructure:"default_lang"`          // 默认语言
	DefaultLocaleKey    string `mapstructure:"default_locale_key"`    // 请求语言key
	MaxMultipartMemory  int64  `mapstructure:"max_multipart_memory"`  // 上传文件最大字节数
	MaxUploadLength     int    `mapstructure:"max_upload_length"`     // 批量上传时每次最多上传多少个文件
	UploadAttachmentDir string `mapstructure:"upload_attachment_dir"` // 前端上传文件的目录
	EmailTemplateDir    string `mapstructure:"email_template_dir"`    // 邮件模板目录
	TemplateDir         string `mapstructure:"template_dir"`          // 模板目录
	Admin               Admin  `mapstructure:"admin"`                 // 管理员信息
	Email               Email  `mapstructure:"email"`                 // 邮件配置
	Hooks               Hooks  `mapstructure:"hooks"`                 // hooks
}

type Admin struct {
	Email string `mapstructure:"email"` // 邮箱
}
type Email struct {
	Host     string `mapstructure:"host"`     // host
	Port     int    `mapstructure:"port"`     // port
	User     string `mapstructure:"host"`     // user
	Password string `mapstructure:"password"` // password
}

type Hooks struct {
	Launch []string `mapstructure:"launch"`
}
