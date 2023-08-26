package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var Module = fx.Options(
	fx.Provide(func() trace.Tracer {
		return otel.Tracer("go-template")
	}),

	fx.Provide(func() (*jaeger.Exporter, error) {
		return jaeger.New(jaeger.WithCollectorEndpoint())
	}),

	fx.Provide(func(exporter *jaeger.Exporter) *sdktrace.TracerProvider {
		provider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("go-template"),
			)),
		)

		otel.SetTracerProvider(provider)

		return provider
	}),

	fx.Invoke(func(livecycle fx.Lifecycle, provider *sdktrace.TracerProvider) {
		livecycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return provider.Shutdown(context.Background())
			},
		})
	}),
)
