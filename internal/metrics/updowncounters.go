package metrics

import (
	"context"

	"github.com/david/otel-datagen/internal/randomness"
	"go.opentelemetry.io/otel/metric"
)

// GenerateInt64UpDownCounter generates int64 updowncounter metrics
func GenerateInt64UpDownCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	upDownCounter, err := meter.Int64UpDownCounter(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		// UpDownCounters can go negative, so we'll allow negative values
		value := int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		if randomness.Float64() < 0.3 { // 30% chance of negative value
			value = -value
		}
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		upDownCounter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}

// GenerateFloat64UpDownCounter generates float64 updowncounter metrics
func GenerateFloat64UpDownCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	upDownCounter, err := meter.Float64UpDownCounter(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		// UpDownCounters can go negative, so we'll allow negative values
		value := float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		if randomness.Float64() < 0.3 { // 30% chance of negative value
			value = -value
		}
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		upDownCounter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}