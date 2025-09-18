package initialize

import (
	"github.com/supuwoerc/weaver/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewTracerProvider(conf *conf.Config, exporter tracesdk.SpanExporter) *tracesdk.TracerProvider {
	tp := tracesdk.NewTracerProvider(
		// 将导出器添加到TracerProvider
		tracesdk.WithBatcher(exporter),
		// 记录关于此应用程序的信息
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.AppInfo()),
			attribute.String("env", conf.Env),
			attribute.String("app_version", conf.AppVersion),
		)),
		// 设置采样率,基于父Span或者配置定义的采样率
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(conf.OpenTelemetry.TraceIDRatioBased))),
	)
	// 设置全局TracerProvider
	otel.SetTracerProvider(tp)
	// 设置传播器,用于跨服务传播追踪上下文
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp
}
