package conf

type OpenTelemetryConfig struct {
	TraceIDRatioBased float64 `mapstructure:"trace_id_ratio_based"` // 采样比例
}
