package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var Module = fx.Options(
	fx.Provide(func() trace.Tracer {
		return otel.Tracer("go-template")
	}),

	fx.Provide(func() *otlptrace.Exporter {
		return otlptracegrpc.NewUnstarted(
			otlptracegrpc.WithEndpoint("localhost:4317"),
			otlptracegrpc.WithInsecure(),
		)
	}),

	fx.Provide(func(exporter *otlptrace.Exporter) *sdktrace.TracerProvider {
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

	fx.Invoke(func(livecycle fx.Lifecycle, exporter *otlptrace.Exporter, provider *sdktrace.TracerProvider) {
		livecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return exporter.Start(ctx)
			},
			OnStop: func(context.Context) error {
				return provider.Shutdown(context.Background())
			},
		})
	}),
)
