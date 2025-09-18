package initialize

import (
	"context"

	"github.com/supuwoerc/weaver/conf"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewOLTPExporter(conf *conf.Config) *otlptrace.Exporter {
	if conf.OLTP.Endpoint == "" {
		panic("oltp endpoint is empty")
	}
	// 使用 OTLP HTTP 导出器
	options := []otlptracehttp.Option{otlptracehttp.WithEndpoint(conf.OLTP.Endpoint)}
	if conf.OLTP.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}
	exp, err := otlptracehttp.New(context.Background(), options...)
	if err != nil {
		panic(err)
	}
	return exp
}
