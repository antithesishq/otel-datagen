package metrics

import (
	"context"

	"github.com/antithesishq/otel-datagen/internal/randomness"
	"go.opentelemetry.io/otel/metric"
)

// GenerateInt64Gauge generates int64 gauge metrics
func GenerateInt64Gauge(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	gauge, err := meter.Int64Gauge(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		gauge.Record(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}

// GenerateFloat64Gauge generates float64 gauge metrics
func GenerateFloat64Gauge(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	gauge, err := meter.Float64Gauge(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		gauge.Record(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}
