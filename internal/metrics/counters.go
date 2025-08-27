package metrics

import (
	"context"

	"github.com/antithesishq/otel-datagen/internal/randomness"
	"github.com/go-faker/faker/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// generateMetricAttributes creates attribute slice for metrics with aggro condition support
func generateMetricAttributes(i int, aggroValues []string, aggroProb float64) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("fake.attr", faker.Word()))
	attrs = append(attrs, attribute.Int("iteration", i+1))

	// Add aggro value attribute if probability check passes
	if aggroProb > 0 && len(aggroValues) > 0 && randomness.Float64() < aggroProb {
		aggroValue := randomness.Choice(aggroValues)
		attrs = append(attrs, attribute.String("aggro.value", aggroValue))
	}

	return attrs
}

// GenerateInt64Counter generates int64 counter metrics
func GenerateInt64Counter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	counter, err := meter.Int64Counter(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		counter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}

// GenerateFloat64Counter generates float64 counter metrics
func GenerateFloat64Counter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	counter, err := meter.Float64Counter(metricName)
	if err != nil {
		return err
	}
	for i := 0; i < numMetrics; i++ {
		value := float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attrs := generateMetricAttributes(i, aggroValues, aggroProb)
		counter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
	return nil
}
