package conf

import (
	"fmt"
	"strings"

	"github.com/supuwoerc/weaver/pkg/constant"
)

// OSSConfig 存储服务配置
type OSSConfig struct {
	Type string `mapstructure:"type"`

	Endpoint string `mapstructure:"endpoint"`
	Region   string `mapstructure:"region"`

	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	ForcePathStyle  bool   `mapstructure:"force_path_style"`

	// 连接配置
	MaxRetries int `mapstructure:"max_retries"`
}

// Validate 验证OSS配置
func (c *OSSConfig) Validate() error {
	if strings.TrimSpace(c.Type) == "" {
		return fmt.Errorf("storage type is required")
	}
	if c.AccessKeyID == "" {
		return fmt.Errorf("access key ID is required")
	}
	if c.SecretAccessKey == "" {
		return fmt.Errorf("secret access key is required")
	}
	if c.Type != constant.AWSS3 && c.Endpoint == "" {
		return fmt.Errorf("endpoint is required for non-AWS S3 services")
	}
	return nil
}
