package exporters

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// CreateResource creates a resource with the given attributes
func CreateResource(ctx context.Context, resourceAttrs []string) (*resource.Resource, error) {
	var attrs []attribute.KeyValue
	
	// Add default service name
	attrs = append(attrs, semconv.ServiceName("otel-datagen"))
	
	// Add custom resource attributes
	for _, attr := range resourceAttrs {
		parts := strings.SplitN(attr, "=", 2)
		if len(parts) == 2 {
			attrs = append(attrs, attribute.String(parts[0], parts[1]))
		}
	}
	
	return resource.New(ctx, resource.WithAttributes(attrs...))
}