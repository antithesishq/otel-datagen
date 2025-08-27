package generators

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/david/otel-datagen/internal/aggro"
	"github.com/david/otel-datagen/internal/exporters"
	"github.com/david/otel-datagen/internal/randomness"
	"github.com/david/otel-datagen/internal/timestamps"
	"github.com/go-faker/faker/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// GenerateTraces generates trace data with the given parameters
func GenerateTraces(numTraces int, numSpans int, numAttributes int, overrideAttrs []string, resourceAttrs []string, otlpEndpoint string, stdoutEnabled bool, timestampConfig *timestamps.TimestampConfig) {
	// Parse aggro configuration for this component
	aggroConfig := aggro.ParseAggroConfig("traces")
	
	ctx := context.Background()

	// Create exporter configuration
	exporterConfig := exporters.ExporterConfig{
		OTLPEndpoint:  otlpEndpoint,
		Protocol:      "grpc", // Default to gRPC for traces
		Insecure:      true,   // For local testing, can be made configurable
		StdoutEnabled: stdoutEnabled,
	}

	// Create dual trace exporters (console + OTLP when endpoint specified)
	traceExporters, err := exporters.CreateDualTraceExporters(ctx, exporterConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create trace exporters: %v", err)
	}

	// Create resource
	res, err := exporters.CreateResource(ctx, resourceAttrs)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Create tracer provider with multiple processors (one per exporter)
	var spanProcessors []trace.TracerProviderOption
	for _, exporter := range traceExporters {
		spanProcessors = append(spanProcessors, trace.WithBatcher(exporter))
	}
	spanProcessors = append(spanProcessors, trace.WithResource(res))

	tp := trace.NewTracerProvider(spanProcessors...)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Generate traces with the provider
	if err := GenerateTracesWithProvider(ctx, tp, numTraces, numSpans, numAttributes, overrideAttrs, aggroConfig, timestampConfig); err != nil {
		log.Printf("Error generating traces: %v", err)
	}

	// Force flush to ensure output
	if err := tp.ForceFlush(ctx); err != nil {
		log.Printf("Error flushing tracer: %v", err)
	}
}

// GenerateTracesWithProvider generates traces using the provided tracer provider
func GenerateTracesWithProvider(ctx context.Context, tp *trace.TracerProvider, numTraces int, numSpans int, numAttributes int, overrideAttrs []string, aggroConfig *aggro.AggroConfig, timestampConfig *timestamps.TimestampConfig) error {
	// Get tracer
	tracer := otel.Tracer("otel-datagen")

	// Parse override attributes
	overrides := make(map[string]string)
	for _, attr := range overrideAttrs {
		parts := strings.SplitN(attr, "=", 2)
		if len(parts) == 2 {
			overrides[parts[0]] = parts[1]
		}
	}


	// Generate the specified number of traces
	totalSpans := 0
	for traceIdx := 0; traceIdx < numTraces; traceIdx++ {
		// Calculate timestamp for the first span of this trace (will be used for root span)
		rootStartTime := timestampConfig.CalculateTimestamp(totalSpans)
		
		// Create a new trace context for this trace with the calculated timestamp
		traceCtx, rootSpan := tracer.Start(ctx, fmt.Sprintf("trace-%d-root", traceIdx+1), oteltrace.WithTimestamp(rootStartTime))
		
		// Generate the specified number of spans for this trace
		for spanIdx := 0; spanIdx < numSpans; spanIdx++ {
			// Calculate timestamp for this span across all traces and spans
			startTime := timestampConfig.CalculateTimestamp(totalSpans)
			totalSpans++

			var span oteltrace.Span
			var spanName string
			
			// First span of each trace is the root span we already created
			if spanIdx == 0 {
				span = rootSpan
				spanName = fmt.Sprintf("trace-%d-root", traceIdx+1)
				// Root span already has the correct start time from tracer.Start()
			} else {
				// Child spans share the same trace context
				spanName = fmt.Sprintf("trace-%d-span-%d", traceIdx+1, spanIdx+1)
				_, span = tracer.Start(traceCtx, spanName, oteltrace.WithTimestamp(startTime))
			}

			// Create attributes list starting with base attribute
			var attrs []attribute.KeyValue
			attrs = append(attrs, semconv.HTTPMethodKey.String("GET"))

			// Generate random attributes using faker
			for j := 0; j < numAttributes; j++ {
				key := fmt.Sprintf("fake.attr.%d", j+1)
				value := faker.Word()

				// Check for override
				if override, exists := overrides[key]; exists {
					value = override
				}

				attrs = append(attrs, attribute.String(key, value))
			}

			// Apply any remaining overrides that didn't match generated keys
			for key, value := range overrides {
				if !strings.HasPrefix(key, "fake.attr.") {
					attrs = append(attrs, attribute.String(key, value))
				}
			}

			// Apply aggro modifications if configured
			skipKeys := []string{"http.method"} // System attributes that shouldn't be replaced
			modifiedAttrs, metadataAttrs := aggroConfig.ApplyAggroToTraceAttributes(attrs, skipKeys, "grpc")
			attrs = modifiedAttrs
			
			// Add metadata attributes about aggro modifications
			attrs = append(attrs, metadataAttrs...)

			span.SetAttributes(attrs...)

			// End span with a duration of 10-100ms after start time
			spanDuration := time.Duration(randomness.Intn(90)+10) * time.Millisecond
			endTime := startTime.Add(spanDuration)
			span.End(oteltrace.WithTimestamp(endTime))
		}
	}

	return nil
}
