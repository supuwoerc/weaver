package conf

import "time"

type ElasticsearchConfig struct {
	Addresses             []string      `mapstructure:"addresses"`               // addresses
	Username              string        `mapstructure:"username"`                // 用户
	Password              string        `mapstructure:"password"`                // password
	Insecure              bool          `mapstructure:"insecure"`                // 关闭TLS
	MaxRetries            int           `mapstructure:"max_retries"`             // 最大重试次数
	CompressRequestBody   bool          `mapstructure:"compress_request_body"`   // 启用压缩
	EnableMetrics         bool          `mapstructure:"enable_metrics"`          // 启用调试metrics
	EnableDebugLogger     bool          `mapstructure:"enable_debug_logger"`     // 启用调试logger
	DiscoverNodesOnStart  bool          `mapstructure:"discover_nodes_on_start"` // 启用节点发现
	DiscoverNodesInterval time.Duration `mapstructure:"discover_nodes_interval"` // 节点发现间隔
	LogLevel              int           `mapstructure:"log_level"`               // 日志级别
}
