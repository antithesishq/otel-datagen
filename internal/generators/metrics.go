package generators

import (
	"context"
	"log"
	"time"

	"github.com/david/otel-datagen/internal/aggro"
	"github.com/david/otel-datagen/internal/exporters"
	"github.com/david/otel-datagen/internal/metrics"
	"github.com/david/otel-datagen/internal/randomness"
	"github.com/david/otel-datagen/internal/timestamps"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// GenerateMetrics generates metric data with the given parameters
func GenerateMetrics(numMetrics int, metricType string, metricName string, counterMin int, counterMax int, resourceAttrs []string, otlpEndpoint string, otlpProtocol string, stdoutEnabled bool, timestampConfig *timestamps.TimestampConfig) {
	// Parse aggro configuration for this component
	aggroConfig := aggro.ParseAggroConfig("metrics")
	
	ctx := context.Background()

	// Create exporter configuration
	exporterConfig := exporters.ExporterConfig{
		OTLPEndpoint:  otlpEndpoint,
		Protocol:      otlpProtocol, // "grpc" or "http"
		Insecure:      true,         // For local testing, can be made configurable
		StdoutEnabled: stdoutEnabled,
	}

	// Create dual metric exporters (console + OTLP when endpoint specified)
	metricExporters, err := exporters.CreateDualMetricExporters(ctx, exporterConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create metric exporters: %v", err)
	}

	// Create resource
	res, err := exporters.CreateResource(ctx, resourceAttrs)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Check if we need timestamp control or regular periodic collection
	if timestampConfig.Spacing > 0 {
		// Use manual readers for timestamp control - need one manual reader for collection
		// but separate exporters for console and OTLP
		reader := sdkmetric.NewManualReader()
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(reader),
			sdkmetric.WithResource(res),
		)
		defer func() {
			if err := mp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down meter provider: %v", err)
			}
		}()

		// Generate metrics with timestamp control using all exporters
		if err := GenerateMetricsWithTimestamps(ctx, mp, reader, metricExporters, numMetrics, metricType, metricName, counterMin, counterMax, timestampConfig); err != nil {
			log.Printf("Error generating metrics: %v", err)
		}
	} else {
		// Use periodic readers for regular operation - one per exporter
		var options []sdkmetric.Option
		for _, exporter := range metricExporters {
			options = append(options, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
		}
		options = append(options, sdkmetric.WithResource(res))
		
		mp := sdkmetric.NewMeterProvider(options...)
		defer func() {
			if err := mp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down meter provider: %v", err)
			}
		}()

		// Generate metrics with the provider
		if err := GenerateMetricsWithProvider(ctx, mp, numMetrics, metricType, metricName, counterMin, counterMax, aggroConfig); err != nil {
			log.Printf("Error generating metrics: %v", err)
		}

		// Force flush to ensure output
		if err := mp.ForceFlush(ctx); err != nil {
			log.Printf("Error flushing meter: %v", err)
		}
	}
}

// GenerateMetricsWithProvider generates metrics using the provided meter provider
func GenerateMetricsWithProvider(ctx context.Context, mp *sdkmetric.MeterProvider, numMetrics int, metricType string, metricName string, counterMin int, counterMax int, aggroConfig *aggro.AggroConfig) error {
	// For now, always use gRPC sanitization in metrics (will be made conditional later)
	return metrics.Generate(ctx, mp, numMetrics, metricType, metricName, counterMin, counterMax, aggroConfig, "grpc")
}

// GenerateMetricsWithTimestamps generates metrics with explicit timestamps using manual reader
func GenerateMetricsWithTimestamps(ctx context.Context, mp *sdkmetric.MeterProvider, reader *sdkmetric.ManualReader, exporters []sdkmetric.Exporter, numMetrics int, metricType string, metricName string, counterMin int, counterMax int, timestampConfig *timestamps.TimestampConfig) error {
	// Generate a time series by creating individual data points at different timestamps
	// Each data point gets fresh random values to create natural variation
	
	for i := 0; i < numMetrics; i++ {
		// Create individual meter provider for this timestamp
		individualReader := sdkmetric.NewManualReader()
		
		// Get resource from original meter provider (reuse existing resource)
		resourceMetrics := &metricdata.ResourceMetrics{}
		if err := reader.Collect(ctx, resourceMetrics); err != nil {
			return err
		}
		
		// Create new meter provider with same resource
		individualMP := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(individualReader),
			sdkmetric.WithResource(resourceMetrics.Resource),
		)
		defer func(mp *sdkmetric.MeterProvider) {
			mp.Shutdown(ctx)
		}(individualMP)
		
		// Generate metrics with timestamped values - this creates fresh random values each time
		if err := generateSingleTimestampedMetric(ctx, individualMP, metricType, metricName, counterMin, counterMax, i); err != nil {
			return err
		}
		
		// Calculate the intended timestamp for this data point
		intendedTimestamp := timestampConfig.CalculateTimestamp(i)
		
		// Collect metrics from the individual provider
		individualResourceMetrics := &metricdata.ResourceMetrics{}
		if err := individualReader.Collect(ctx, individualResourceMetrics); err != nil {
			return err
		}
		
		// Manually adjust timestamps in the collected data to match intended timestamp
		adjustTimestamps(individualResourceMetrics, intendedTimestamp)
		
		// Export the timestamped metrics to all exporters
		for _, exporter := range exporters {
			if err := exporter.Export(ctx, individualResourceMetrics); err != nil {
				return err
			}
		}
		
		// Small delay between iterations for cleaner output
		if i < numMetrics-1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	
	return nil
}

// generateSingleTimestampedMetric generates a single metric data point with fresh random values
func generateSingleTimestampedMetric(ctx context.Context, mp *sdkmetric.MeterProvider, metricType string, metricName string, counterMin int, counterMax int, iteration int) error {
	// Get meter
	meter := mp.Meter("otel-datagen")
	
	// Generate fresh attributes for this iteration
	attrs := generateTimestampedMetricAttributes(iteration)
	
	// Generate fresh random value for this timestamp
	value := generateFreshRandomValue(counterMin, counterMax)
	
	// Create the metric based on type and record the value
	switch metricType {
	case "gauge", "int64-gauge":
		gauge, err := meter.Int64Gauge(metricName)
		if err != nil {
			return err
		}
		gauge.Record(ctx, int64(value), metric.WithAttributes(attrs...))
		
	case "float64-gauge":
		gauge, err := meter.Float64Gauge(metricName)
		if err != nil {
			return err
		}
		gauge.Record(ctx, float64(value), metric.WithAttributes(attrs...))
		
	case "counter", "int64-counter":
		counter, err := meter.Int64Counter(metricName)
		if err != nil {
			return err
		}
		counter.Add(ctx, int64(value), metric.WithAttributes(attrs...))
		
	case "float64-counter":
		counter, err := meter.Float64Counter(metricName)
		if err != nil {
			return err
		}
		counter.Add(ctx, float64(value), metric.WithAttributes(attrs...))
		
	case "histogram", "float64-histogram":
		histogram, err := meter.Float64Histogram(metricName)
		if err != nil {
			return err
		}
		histogram.Record(ctx, float64(value), metric.WithAttributes(attrs...))
		
	case "int64-histogram":
		histogram, err := meter.Int64Histogram(metricName)
		if err != nil {
			return err
		}
		histogram.Record(ctx, int64(value), metric.WithAttributes(attrs...))
		
	default:
		// For other metric types, fall back to the original generation method
		return metrics.Generate(ctx, mp, 1, metricType, metricName, counterMin, counterMax, nil, "grpc")
	}
	
	return nil
}

// generateTimestampedMetricAttributes creates consistent attributes for time-series data
func generateTimestampedMetricAttributes(iteration int) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	
	// Use absolutely minimal consistent attributes to ensure single time series
	// No changing attributes that would create separate series in Prometheus
	attrs = append(attrs, attribute.String("instance", "primary"))
	
	// Don't add conditional attributes - aggro probability will be handled in value generation
	// This ensures all data points belong to the same time series
	
	return attrs
}

// generateFreshRandomValue creates a new random value for each call, ensuring variation
func generateFreshRandomValue(counterMin, counterMax int) int {
	// Generate normal random value in range
	return randomness.Intn(counterMax-counterMin+1) + counterMin
}

// adjustTimestamps manually adjusts timestamps in metric data to the intended timestamp
// This is a workaround for OpenTelemetry Go SDK's limitation in timestamp control
func adjustTimestamps(resourceMetrics *metricdata.ResourceMetrics, intendedTimestamp time.Time) {
	for i := range resourceMetrics.ScopeMetrics {
		for j := range resourceMetrics.ScopeMetrics[i].Metrics {
			metric := &resourceMetrics.ScopeMetrics[i].Metrics[j]
			
			switch data := metric.Data.(type) {
			case metricdata.Gauge[int64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
					// For gauges, StartTime is typically the same as Time
					data.DataPoints[k].StartTime = intendedTimestamp
				}
			case metricdata.Gauge[float64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
					data.DataPoints[k].StartTime = intendedTimestamp
				}
			case metricdata.Sum[int64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
					// For cumulative sums, keep original StartTime but update Time
					if !data.IsMonotonic {
						data.DataPoints[k].StartTime = intendedTimestamp
					}
				}
			case metricdata.Sum[float64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
					if !data.IsMonotonic {
						data.DataPoints[k].StartTime = intendedTimestamp
					}
				}
			case metricdata.Histogram[int64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
				}
			case metricdata.Histogram[float64]:
				for k := range data.DataPoints {
					data.DataPoints[k].Time = intendedTimestamp
				}
			}
		}
	}
}