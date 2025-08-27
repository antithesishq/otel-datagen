package metrics

import (
	"context"
	"fmt"

	"github.com/antithesishq/otel-datagen/internal/aggro"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Generate generates metrics using the provided meter provider
func Generate(ctx context.Context, mp *sdkmetric.MeterProvider, numMetrics int, metricType string, metricName string, counterMin int, counterMax int, aggroConfig *aggro.AggroConfig, protocol string) error {
	// Get meter
	meter := mp.Meter("otel-datagen")

	// Generate aggro values for aggro system
	var aggroValues []string
	var effectiveAggroProb float64

	if aggroConfig != nil && aggroConfig.HasAnyActive() {
		// Use aggro values when aggro is configured - set probability to 1.0 to always trigger
		aggroValues = aggro.GetAggroValuesForProtocol(protocol)
		effectiveAggroProb = 1.0
	}

	// Generate metrics based on type
	switch metricType {
	case "counter", "int64-counter":
		return GenerateInt64Counter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-counter":
		return GenerateFloat64Counter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "histogram", "float64-histogram":
		return GenerateFloat64Histogram(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "int64-histogram":
		return GenerateInt64Histogram(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "updowncounter", "int64-updowncounter":
		return GenerateInt64UpDownCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-updowncounter":
		return GenerateFloat64UpDownCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "gauge", "int64-gauge":
		return GenerateInt64Gauge(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-gauge":
		return GenerateFloat64Gauge(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "observable-counter", "int64-observable-counter":
		return GenerateInt64ObservableCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-observable-counter":
		return GenerateFloat64ObservableCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "observable-updowncounter", "int64-observable-updowncounter":
		return GenerateInt64ObservableUpDownCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-observable-updowncounter":
		return GenerateFloat64ObservableUpDownCounter(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "observable-gauge", "int64-observable-gauge":
		return GenerateInt64ObservableGauge(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	case "float64-observable-gauge":
		return GenerateFloat64ObservableGauge(ctx, meter, metricName, numMetrics, counterMin, counterMax, aggroValues, effectiveAggroProb)
	default:
		return fmt.Errorf("unsupported metric type: %s (supported: counter, float64-counter, histogram, int64-histogram, updowncounter, float64-updowncounter, gauge, float64-gauge, observable-counter, float64-observable-counter, observable-updowncounter, float64-observable-updowncounter, observable-gauge, float64-observable-gauge)", metricType)
	}
}
