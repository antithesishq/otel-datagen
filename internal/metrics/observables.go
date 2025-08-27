package metrics

import (
	"context"

	"github.com/david/otel-datagen/internal/randomness"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// GenerateInt64ObservableCounter generates int64 observable counter metrics
func GenerateInt64ObservableCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]int64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Int64ObservableCounter(metricName, metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}

// GenerateFloat64ObservableCounter generates float64 observable counter metrics
func GenerateFloat64ObservableCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]float64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Float64ObservableCounter(metricName, metric.WithFloat64Callback(func(_ context.Context, observer metric.Float64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}

// GenerateInt64ObservableUpDownCounter generates int64 observable updowncounter metrics
func GenerateInt64ObservableUpDownCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]int64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		if randomness.Float64() < 0.3 { // 30% chance of negative value
			values[i] = -values[i]
		}
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Int64ObservableUpDownCounter(metricName, metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}

// GenerateFloat64ObservableUpDownCounter generates float64 observable updowncounter metrics
func GenerateFloat64ObservableUpDownCounter(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]float64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		if randomness.Float64() < 0.3 { // 30% chance of negative value
			values[i] = -values[i]
		}
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Float64ObservableUpDownCounter(metricName, metric.WithFloat64Callback(func(_ context.Context, observer metric.Float64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}

// GenerateInt64ObservableGauge generates int64 observable gauge metrics
func GenerateInt64ObservableGauge(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]int64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = int64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Int64ObservableGauge(metricName, metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}

// GenerateFloat64ObservableGauge generates float64 observable gauge metrics
func GenerateFloat64ObservableGauge(ctx context.Context, meter metric.Meter, metricName string, numMetrics int, counterMin int, counterMax int, aggroValues []string, aggroProb float64) error {
	// Generate values to be used in callback
	values := make([]float64, numMetrics)
	attributeSets := make([][]attribute.KeyValue, numMetrics)
	for i := 0; i < numMetrics; i++ {
		values[i] = float64(randomness.Intn(counterMax-counterMin+1) + counterMin)
		attributeSets[i] = generateMetricAttributes(i, aggroValues, aggroProb)
	}
	
	_, err := meter.Float64ObservableGauge(metricName, metric.WithFloat64Callback(func(_ context.Context, observer metric.Float64Observer) error {
		for i, value := range values {
			observer.Observe(value, metric.WithAttributes(attributeSets[i]...))
		}
		return nil
	}))
	return err
}