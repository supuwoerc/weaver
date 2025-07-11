package conf

import "time"

type LoggerConfig struct {
	MaxSize                       int           `mapstructure:"max_size"`                           // 日志文件切割尺寸(m)
	MaxBackups                    int           `mapstructure:"max_backups"`                        // 保留文件对最大个数
	MaxAge                        int           `mapstructure:"max_age"`                            // 保留文件对最大天数
	Level                         int8          `mapstructure:"level"`                              // 日志级别
	Dir                           string        `mapstructure:"dir"`                                // 日志文件存放的目录,为空时默认在项目目录下创建log目录存放日志文件
	Stdout                        bool          `mapstructure:"stdout"`                             // 标准终端输出
	GormLevel                     int           `mapstructure:"gorm_level"`                         // gorm日志级别
	GormSlowThreshold             time.Duration `mapstructure:"gorm_slow_threshold"`                // gorm慢查询阈值
	GormIgnoreRecordNotFoundError bool          `mapstructure:"gorm_ignore_record_not_found_error"` // gorm忽略record not found错误
	RedisLevel                    int           `mapstructure:"redis_level"`                        // redis日志级别
}
