package metrics

import (
	"context"

	"github.com/antithesishq/otel-datagen/internal/randomness"
	"go.opentelemetry.io/otel/metric"
)

// GenerateInt64Histogram generates int64 histogram metrics
func GenerateInt64Histogram(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	histogram, err := meter.Int64Histogram(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		histogram.Record(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}

// GenerateFloat64Histogram generates float64 histogram metrics
func GenerateFloat64Histogram(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	histogram, err := meter.Float64Histogram(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		histogram.Record(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}
