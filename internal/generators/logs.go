package generators

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/antithesishq/otel-datagen/internal/aggro"
	"github.com/antithesishq/otel-datagen/internal/exporters"
	"github.com/antithesishq/otel-datagen/internal/timestamps"
	"github.com/go-faker/faker/v4"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// GenerateLogs generates log data with the given parameters
func GenerateLogs(numLogs int, numAttributes int, overrideAttrs []string, resourceAttrs []string, otlpEndpoint string, otlpProtocol string, stdoutEnabled bool, timestampConfig *timestamps.TimestampConfig) {
	// Parse aggro configuration for this component
	aggroConfig := aggro.ParseAggroConfig("logs")

	ctx := context.Background()

	// Create exporter configuration
	exporterConfig := exporters.ExporterConfig{
		OTLPEndpoint:  otlpEndpoint,
		Protocol:      otlpProtocol, // "grpc" or "http"
		Insecure:      true,         // For local testing, can be made configurable
		StdoutEnabled: stdoutEnabled,
	}

	// Create dual log exporters (console + OTLP when endpoint specified)
	logExporters, err := exporters.CreateDualLogExporters(ctx, exporterConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create log exporters: %v", err)
	}

	// Create resource
	res, err := exporters.CreateResource(ctx, resourceAttrs)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Create logger provider with multiple processors (one per exporter)
	var logProcessors []sdklog.LoggerProviderOption
	for _, exporter := range logExporters {
		logProcessors = append(logProcessors, sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)))
	}
	logProcessors = append(logProcessors, sdklog.WithResource(res))

	lp := sdklog.NewLoggerProvider(logProcessors...)
	defer func() {
		if err := lp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down logger provider: %v", err)
		}
	}()

	// Generate logs with the provider
	if err := GenerateLogsWithProvider(ctx, lp, numLogs, numAttributes, overrideAttrs, aggroConfig, timestampConfig); err != nil {
		log.Printf("Error generating logs: %v", err)
	}

	// Force flush to ensure output
	if err := lp.ForceFlush(ctx); err != nil {
		log.Printf("Error flushing logger: %v", err)
	}
}

// GenerateLogsWithProvider generates logs using the provided logger provider
func GenerateLogsWithProvider(ctx context.Context, lp *sdklog.LoggerProvider, numLogs int, numAttributes int, overrideAttrs []string, aggroConfig *aggro.AggroConfig, timestampConfig *timestamps.TimestampConfig) error {
	// Get logger
	logger := lp.Logger("otel-datagen")

	// Parse override attributes
	overrides := make(map[string]string)
	for _, attr := range overrideAttrs {
		parts := strings.SplitN(attr, "=", 2)
		if len(parts) == 2 {
			overrides[parts[0]] = parts[1]
		}
	}

	// Generate the specified number of log records
	for i := 0; i < numLogs; i++ {
		logMessage := fmt.Sprintf("example-log-%d", i+1)

		// Create attributes list starting with base attribute
		var attrs []otellog.KeyValue
		attrs = append(attrs, otellog.String("log.level", "info"))

		// Generate random attributes using faker
		for j := 0; j < numAttributes; j++ {
			key := fmt.Sprintf("fake.attr.%d", j+1)
			value := faker.Word()

			// Check for override
			if override, exists := overrides[key]; exists {
				value = override
			}

			attrs = append(attrs, otellog.String(key, value))
		}

		// Apply any remaining overrides that didn't match generated keys
		for key, value := range overrides {
			if !strings.HasPrefix(key, "fake.attr.") {
				attrs = append(attrs, otellog.String(key, value))
			}
		}

		// Apply aggro modifications if configured
		skipKeys := []string{"log.level"} // System attributes that shouldn't be replaced
		// Use gRPC sanitization for now (will be made conditional in next iteration)
		modifiedAttrs, modifiedMessage, metadataAttrs := aggroConfig.ApplyAggroToLogAttributes(attrs, logMessage, skipKeys, "grpc")
		attrs = modifiedAttrs
		logMessage = modifiedMessage

		// Add metadata attributes about aggro modifications
		attrs = append(attrs, metadataAttrs...)

		record := otellog.Record{}
		record.SetBody(otellog.StringValue(logMessage))
		record.SetSeverity(otellog.SeverityInfo)

		// Set timestamp for this log record
		logTime := timestampConfig.CalculateTimestamp(i)
		record.SetTimestamp(logTime)
		// Also set observed timestamp to match for historical data generation
		record.SetObservedTimestamp(logTime)

		record.AddAttributes(attrs...)
		logger.Emit(ctx, record)
	}

	return nil
}
