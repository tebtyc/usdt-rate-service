package tracing

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitTracer инициализирует трассировщик для вашего сервиса
func InitTracer(jaegerURL, serviceName string) (*tracesdk.TracerProvider, error) {
	// Создаем Jaeger экспортёр
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, fmt.Errorf("creating Jaeger exporter: %w", err)
	}

	// Создаем ресурс, который будет содержать информацию о сервисе
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Создаем провайдер трассировки
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	// Устанавливаем провайдер для OpenTelemetry
	otel.SetTracerProvider(tp)

	return tp, nil
}
