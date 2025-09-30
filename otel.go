package main

import (
	"context"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Create OTLP trace exporter over gRPC (default localhost:4317)
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	// Identify your service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(consts.APP_NAME),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create tracer provider with batch span processor
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Register it as global provider
	otel.SetTracerProvider(tp)
	return tp, nil
}
