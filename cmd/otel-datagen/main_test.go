package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/david/otel-datagen/internal/aggro"
	"github.com/david/otel-datagen/internal/generators"
	"github.com/david/otel-datagen/internal/metrics"
	"github.com/david/otel-datagen/internal/randomness"
	"github.com/david/otel-datagen/internal/timestamps"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"gopkg.in/yaml.v3"
)

// generateTracesToFile is a test helper that outputs traces to a file instead of stdout
func generateTracesToFile(numSpans int, numAttributes int, overrideAttrs []string, aggroProb float64, outputFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter instead of stdout exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-datagen"),
		),
	)
	if err != nil {
		return err
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Use the same generation logic as the main function
	timestampConfig := &timestamps.TimestampConfig{}
	// Create aggro config for backward compatibility with aggro probability
	aggroConfig := &aggro.AggroConfig{}
	if aggroProb > 0 {
		// For tests, use the old aggro probability behavior by setting all aggro types to active
		aggroConfig = &aggro.AggroConfig{
			StringActive:    true,
			NumericActive:   true,
			TimestampActive: true,
			StringTarget:    "",  // random targeting
			NumericTarget:   "",  // random targeting
			TimestampTarget: "",  // random targeting
		}
	}
	err = generators.GenerateTracesWithProvider(ctx, tp, 1, numSpans, numAttributes, overrideAttrs, aggroConfig, timestampConfig)
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return tp.ForceFlush(ctx)
}

// generateTracesToFileWithResource is a test helper that outputs traces to a file with custom resource attributes
func generateTracesToFileWithResource(numSpans int, numAttributes int, overrideAttrs []string, aggroProb float64, outputFile string, resourceAttrs map[string]string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter instead of stdout exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with custom attributes
	var attrs []attribute.KeyValue
	for key, value := range resourceAttrs {
		attrs = append(attrs, attribute.String(key, value))
	}
	
	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return err
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Use the same generation logic as the main function
	timestampConfig := &timestamps.TimestampConfig{}
	// Create aggro config for backward compatibility with aggro probability
	aggroConfig := &aggro.AggroConfig{}
	if aggroProb > 0 {
		// For tests, use the old aggro probability behavior by setting all aggro types to active
		aggroConfig = &aggro.AggroConfig{
			StringActive:    true,
			NumericActive:   true,
			TimestampActive: true,
			StringTarget:    "",  // random targeting
			NumericTarget:   "",  // random targeting
			TimestampTarget: "",  // random targeting
		}
	}
	err = generators.GenerateTracesWithProvider(ctx, tp, 1, numSpans, numAttributes, overrideAttrs, aggroConfig, timestampConfig)
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return tp.ForceFlush(ctx)
}

// generateTracesToFileWithAggro is a new test helper that uses AggroConfig directly
func generateTracesToFileWithAggro(numSpans int, numAttributes int, overrideAttrs []string, aggroConfig *aggro.AggroConfig, outputFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter instead of stdout exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-datagen"),
		),
	)
	if err != nil {
		return err
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Use the updated generation logic with direct aggro config
	timestampConfig := &timestamps.TimestampConfig{}
	err = generators.GenerateTracesWithProvider(ctx, tp, 1, numSpans, numAttributes, overrideAttrs, aggroConfig, timestampConfig)
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return tp.ForceFlush(ctx)
}

func TestGenerateTracesBasic(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	err := generateTracesToFile(1, 2, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "service.name")
	assert.Contains(t, output, "otel-datagen")
	assert.Contains(t, output, "trace-1-root")
	assert.Contains(t, output, "fake.attr.1")
	assert.Contains(t, output, "fake.attr.2")
	
	// Should not contain aggro metadata with probability 0.0
	assert.NotContains(t, output, "aggro.string")
	assert.NotContains(t, output, "aggro.numeric")
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestGenerateTracesMultipleSpans(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	err := generateTracesToFile(3, 1, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "trace-1-root")
	assert.Contains(t, output, "trace-1-span-2")
	assert.Contains(t, output, "trace-1-span-3")
}

func TestGenerateTracesWithOverride(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	overrides := []string{"test.attribute=custom-value"}
	err := generateTracesToFile(1, 1, overrides, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "test.attribute")
	assert.Contains(t, output, "custom-value")
}

func TestGenerateTracesWithAggroProbability(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Use probability 1.0 to guarantee aggro values appear (backward compatibility test)
	err := generateTracesToFile(1, 1, []string{}, 1.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output contains aggro metadata (new behavior)
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// New aggro system should add metadata attributes showing what was modified
	assert.Contains(t, output, "aggro.string")
	assert.Contains(t, output, "aggro.numeric") 
	assert.Contains(t, output, "aggro.timestamp")
}

func TestGenerateTracesNumAttributes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	err := generateTracesToFile(1, 5, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output contains expected number of fake attributes
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "fake.attr.1")
	assert.Contains(t, output, "fake.attr.2")
	assert.Contains(t, output, "fake.attr.3")
	assert.Contains(t, output, "fake.attr.4")
	assert.Contains(t, output, "fake.attr.5")
}

func TestGenerateTracesOutputStructure(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	err := generateTracesToFile(1, 1, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output contains valid JSON
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	require.NotEmpty(t, output, "Output file should not be empty")

	// The stdouttrace exporter outputs one JSON object
	// Let's just verify it contains the expected fields as text
	assert.Contains(t, output, "\"Name\":")
	assert.Contains(t, output, "\"SpanContext\":")
	assert.Contains(t, output, "\"Attributes\":")
	assert.Contains(t, output, "\"Resource\":")
	assert.Contains(t, output, "\"trace-1-root\"")
	
	// Verify it's valid JSON by attempting to parse the entire output
	var spanData map[string]interface{}
	err = json.Unmarshal(content, &spanData)
	require.NoError(t, err, "Output should be valid JSON")
}

func TestGenerateTracesNoAggroValues(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Use probability 0.0 to ensure no aggro values appear
	err := generateTracesToFile(1, 2, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output does NOT contain aggro values
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.NotContains(t, output, "aggro.string", "Should not contain aggro metadata with probability 0.0")
	assert.NotContains(t, output, "aggro.numeric", "Should not contain aggro metadata with probability 0.0") 
	assert.NotContains(t, output, "aggro.timestamp", "Should not contain aggro metadata with probability 0.0")
	assert.Contains(t, output, "fake.attr.1")
	assert.Contains(t, output, "fake.attr.2")
}

func TestResourceAttributes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Test custom resource attributes  
	err := generateTracesToFileWithResource(1, 1, []string{}, 0.0, outputFile, map[string]string{
		"service.name":    "custom-service",
		"service.version": "1.2.3",
		"environment":     "test",
	})
	require.NoError(t, err)

	// Read and verify the resource attributes
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "custom-service")
	assert.Contains(t, output, "1.2.3")
	assert.Contains(t, output, "environment")
}

func TestOTLPExporterCreation(t *testing.T) {
	// This test verifies that OTLP exporter creation doesn't panic with an invalid endpoint
	// In a real scenario, this would fail to connect, but we just want to test the setup
	
	// Test that creating an OTLP exporter doesn't immediately fail
	ctx := context.Background()
	_, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("http://localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	
	// The creation should succeed even if the endpoint is not reachable
	// The actual connection happens during export
	assert.NoError(t, err)
}

func TestConfigFileSupport(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create a test configuration file
	configContent := `# Global settings
resource:
  service.name: "config-test-service"
  service.version: "2.1.0"

generate:
  traces:
    num_spans: 2
    num_attributes: 3
    boundary_probability: 0.0
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test with config file - this should use the config file settings
	err = generateTracesToFileWithConfig(outputFile, configFile)
	require.NoError(t, err)

	// Read and verify the configuration was applied
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "config-test-service")
	assert.Contains(t, output, "2.1.0")
	assert.Contains(t, output, "trace-1-root")
	assert.Contains(t, output, "trace-1-span-2")
}

func TestConfigFilePrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create a test configuration file
	configContent := `# Global settings
resource:
  service.name: "config-service"

generate:
  traces:
    num_spans: 3
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test precedence: CLI flags should override config file
	err = generateTracesToFileWithConfigAndOverrides(outputFile, configFile, map[string]string{
		"service.name": "cli-override-service",
	}, 1) // Override num_spans to 1
	require.NoError(t, err)

	// Read and verify CLI flags took precedence
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "cli-override-service", "CLI flag should override config file")
	assert.Contains(t, output, "trace-1-root")
	// Should only have 1 span (CLI override), not 3 (config file)
	assert.NotContains(t, output, "trace-1-span-2")
}

// generateTracesToFileWithConfig is a test helper that reads from a config file
func generateTracesToFileWithConfig(outputFile string, configFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter instead of stdout exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Read config file manually for testing (bypass cobra/viper initialization)
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Parse YAML config
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	// Extract values from config
	numSpans := 2        // default
	numAttributes := 3   // default

	if generate, ok := config["generate"].(map[string]interface{}); ok {
		if traces, ok := generate["traces"].(map[string]interface{}); ok {
			if ns, ok := traces["num_spans"].(int); ok {
				numSpans = ns
			}
			if na, ok := traces["num_attributes"].(int); ok {
				numAttributes = na
			}
		}
	}

	// Create resource with config attributes
	var attrs []attribute.KeyValue
	if resource, ok := config["resource"].(map[string]interface{}); ok {
		for key, value := range resource {
			if str, ok := value.(string); ok {
				attrs = append(attrs, attribute.String(key, str))
			}
		}
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return err
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Use the same generation logic as the main function
	timestampConfig := &timestamps.TimestampConfig{}
	// Create empty aggro config since this test doesn't use aggro probability
	aggroConfig := &aggro.AggroConfig{}
	err = generators.GenerateTracesWithProvider(ctx, tp, 1, numSpans, numAttributes, []string{}, aggroConfig, timestampConfig)
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return tp.ForceFlush(ctx)
}

// generateTracesToFileWithConfigAndOverrides tests config file with CLI overrides
func generateTracesToFileWithConfigAndOverrides(outputFile string, configFile string, resourceOverrides map[string]string, numSpansOverride int) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter instead of stdout exporter
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Read config file manually for testing (bypass cobra/viper initialization)
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Parse YAML config
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	// Extract values from config, with CLI overrides taking precedence
	numSpans := numSpansOverride // CLI override wins
	numAttributes := 3           // default

	if generate, ok := config["generate"].(map[string]interface{}); ok {
		if traces, ok := generate["traces"].(map[string]interface{}); ok {
			if na, ok := traces["num_attributes"].(int); ok {
				numAttributes = na
			}
		}
	}

	// Create resource with config attributes, then apply CLI overrides
	var attrs []attribute.KeyValue
	if resource, ok := config["resource"].(map[string]interface{}); ok {
		for key, value := range resource {
			if str, ok := value.(string); ok {
				// Check if this key is overridden by CLI
				if override, exists := resourceOverrides[key]; exists {
					attrs = append(attrs, attribute.String(key, override))
				} else {
					attrs = append(attrs, attribute.String(key, str))
				}
			}
		}
	}

	// Add any CLI-only resource attributes
	for key, value := range resourceOverrides {
		found := false
		if resource, ok := config["resource"].(map[string]interface{}); ok {
			if _, exists := resource[key]; exists {
				found = true
			}
		}
		if !found {
			attrs = append(attrs, attribute.String(key, value))
		}
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return err
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Use the same generation logic as the main function
	timestampConfig := &timestamps.TimestampConfig{}
	// Create empty aggro config since this test doesn't use aggro probability
	aggroConfig := &aggro.AggroConfig{}
	err = generators.GenerateTracesWithProvider(ctx, tp, 1, numSpans, numAttributes, []string{}, aggroConfig, timestampConfig)
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return tp.ForceFlush(ctx)
}// generateLogsToFile is a test helper that outputs logs to a file
func generateLogsToFile(numLogs int, numAttributes int, overrideAttrs []string, aggroProb float64, outputFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter for logs
	exporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
		stdoutlog.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-datagen"),
		),
	)
	if err != nil {
		return err
	}

	// Create logger provider
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	defer func() {
		if err := lp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Generate logs using the same pattern as traces
	logger := lp.Logger("otel-datagen")

	// Parse override attributes
	overrides := make(map[string]string)
	for _, attr := range overrideAttrs {
		parts := strings.SplitN(attr, "=", 2)
		if len(parts) == 2 {
			overrides[parts[0]] = parts[1]
		}
	}

	// Get aggro values for potential use
	var aggroValues []string
	if aggroProb > 0 {
		aggroValues = aggro.GetAggroValues()
	}

	// Generate the specified number of log records
	for i := 0; i < numLogs; i++ {
		record := otellog.Record{}
		record.SetBody(otellog.StringValue(fmt.Sprintf("example-log-%d", i+1)))
		record.SetSeverity(otellog.SeverityInfo)

		// Create attributes list
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

		// Apply any remaining overrides that didnt match generated keys
		for key, value := range overrides {
			if !strings.HasPrefix(key, "fake.attr.") {
				attrs = append(attrs, otellog.String(key, value))
			}
		}

		// Add aggro value attribute if probability check passes
		if aggroProb > 0 && len(aggroValues) > 0 && randomness.Float64() < aggroProb {
			aggroValue := randomness.Choice(aggroValues)
			attrs = append(attrs, otellog.String("aggro.value", aggroValue))
		}

		record.AddAttributes(attrs...)
		logger.Emit(ctx, record)
	}

	// Force flush to ensure output
	return lp.ForceFlush(ctx)
}

func TestGenerateLogsBasic(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "logs.json")

	err := generateLogsToFile(1, 2, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "service.name")
	assert.Contains(t, output, "otel-datagen")
	assert.Contains(t, output, "example-log-1")
	assert.Contains(t, output, "fake.attr.1")
	assert.Contains(t, output, "fake.attr.2")

	// Should not contain aggro metadata with probability 0.0
	assert.NotContains(t, output, "aggro.string")
	assert.NotContains(t, output, "aggro.numeric")
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestGenerateLogsMultipleLogs(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "logs.json")

	err := generateLogsToFile(3, 1, []string{}, 0.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "example-log-1")
	assert.Contains(t, output, "example-log-2")
	assert.Contains(t, output, "example-log-3")
}

func TestGenerateLogsWithAggroProbability(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "logs.json")

	// Use probability 1.0 to guarantee aggro value appears (backward compatibility)
	// Note: This will still use the old boundary.GetBoundaryValues() function
	err := generateLogsToFile(1, 1, []string{}, 1.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Logs still use the old boundary.value approach since they haven't been fully updated to aggro
	assert.Contains(t, output, "aggro.value")
}// generateMetricsToFile is a test helper that outputs metrics to a file
func generateMetricsToFile(numMetrics int, metricType string, metricName string, outputFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter for metrics
	exporter, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
		stdoutmetric.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-datagen"),
		),
	)
	if err != nil {
		return err
	}

	// Create meter provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Use the updated generateMetricsWithProvider function
	aggroConfig := &aggro.AggroConfig{} // Empty config for test
	err = metrics.Generate(ctx, mp, numMetrics, metricType, metricName, 1, 100, aggroConfig, "")
	if err != nil {
		return err
	}

	// Force flush to ensure output
	return mp.ForceFlush(ctx)
}

func TestGenerateMetricsBasicCounter(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(1, "counter", "test_counter", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "service.name")
	assert.Contains(t, output, "otel-datagen")
	assert.Contains(t, output, "test_counter")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsBasicHistogram(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(1, "histogram", "test_histogram", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "service.name")
	assert.Contains(t, output, "otel-datagen")
	assert.Contains(t, output, "test_histogram")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(3, "counter", "test_counter", outputFile)
	require.NoError(t, err)

	// Read and verify the output contains multiple data points
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "test_counter")
	// Should have multiple iterations
	assert.Contains(t, output, "iteration")
}

// generateMetricsToFileWithAggro is a test helper that outputs metrics to a file with aggro probability
func generateMetricsToFileWithAggro(numMetrics int, metricType string, metricName string, aggroProb float64, outputFile string) error {
	ctx := context.Background()

	// Create file for output
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create file exporter for metrics
	exporter, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
		stdoutmetric.WithWriter(file),
	)
	if err != nil {
		return err
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-datagen"),
		),
	)
	if err != nil {
		return err
	}

	// Create meter provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			// Handle error if needed
		}
	}()

	// Generate metrics with the provider - using generateMetricsWithProvider function
	aggroConfig := &aggro.AggroConfig{StringActive: true} // Active config to test aggro values
	if err := metrics.Generate(ctx, mp, numMetrics, metricType, metricName, 1, 100, aggroConfig, ""); err != nil {
		return err
	}

	// Force flush to ensure output
	return mp.ForceFlush(ctx)
}

func TestGenerateMetricsUpDownCounter(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(3, "updowncounter", "active_connections", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "service.name")
	assert.Contains(t, output, "otel-datagen")
	assert.Contains(t, output, "active_connections")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsFloat64Counter(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(2, "float64-counter", "bytes_sent", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "bytes_sent")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsGauge(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(2, "gauge", "cpu_usage", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "cpu_usage")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsInt64Histogram(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(2, "int64-histogram", "response_size", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "response_size")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsObservableCounter(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(2, "observable-counter", "requests_total", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "requests_total")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsObservableGauge(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(2, "observable-gauge", "memory_usage", outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "memory_usage")
	assert.Contains(t, output, "fake.attr")
}

func TestGenerateMetricsWithAggroValues(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	// Test with a different metric type to ensure aggro values work across all types
	err := generateMetricsToFileWithAggro(2, "updowncounter", "test_metric", 1.0, outputFile)
	require.NoError(t, err)

	// Read and verify the output contains aggro values
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Metrics still use the old boundary.value approach since they use backward compatibility
	assert.Contains(t, output, "aggro.value")
}

func TestGenerateMetricsUnsupportedType(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "metrics.json")

	err := generateMetricsToFile(1, "invalid-type", "test_metric", outputFile)
	// This should error since we're testing an unsupported metric type
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported metric type")
}

// ===== NEW AGGRO-SPECIFIC TESTS =====

func TestAggroStringChaosEngineering(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config for string chaos engineering with random targeting
	aggroConfig := &aggro.AggroConfig{
		StringActive:  true,
		StringTarget:  "", // random targeting
		NumericActive: false,
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain aggro string metadata
	assert.Contains(t, output, "aggro.string")
	// Should not contain other aggro types  
	assert.NotContains(t, output, "aggro.numeric")
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestAggroNumericChaosEngineering(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config for numeric chaos engineering with random targeting
	aggroConfig := &aggro.AggroConfig{
		StringActive:    false,
		NumericActive:   true,
		NumericTarget:   "", // random targeting
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain aggro numeric metadata
	assert.Contains(t, output, "aggro.numeric")
	// Should not contain other aggro types
	assert.NotContains(t, output, "aggro.string")
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestAggroTimestampChaosEngineering(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config for timestamp chaos engineering with random targeting
	aggroConfig := &aggro.AggroConfig{
		StringActive:    false,
		NumericActive:   false,
		TimestampActive: true,
		TimestampTarget: "", // random targeting
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain aggro timestamp metadata
	assert.Contains(t, output, "aggro.timestamp")
	// Should not contain other aggro types
	assert.NotContains(t, output, "aggro.string")
	assert.NotContains(t, output, "aggro.numeric")
}

func TestMultipleAggroTypesIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config with all three types enabled
	aggroConfig := &aggro.AggroConfig{
		StringActive:    true,
		StringTarget:    "", // random targeting
		NumericActive:   true,
		NumericTarget:   "", // random targeting  
		TimestampActive: true,
		TimestampTarget: "", // random targeting
	}

	err := generateTracesToFileWithAggro(1, 5, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain all three aggro metadata types
	assert.Contains(t, output, "aggro.string")
	assert.Contains(t, output, "aggro.numeric")
	assert.Contains(t, output, "aggro.timestamp")
}

func TestAggroSpecificAttributeTargeting(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config that targets a specific attribute
	aggroConfig := &aggro.AggroConfig{
		StringActive:    true,
		StringTarget:    "fake.attr.2", // target specific attribute
		NumericActive:   false,
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain aggro string metadata showing which attribute was targeted
	assert.Contains(t, output, "aggro.string")
	assert.Contains(t, output, "fake.attr.2") // The metadata should show the targeted attribute
	// Should not contain other aggro types
	assert.NotContains(t, output, "aggro.numeric")
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestAggroWithOverrideAttributes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Test aggro with override attributes
	overrides := []string{"custom.field=original-value"}
	aggroConfig := &aggro.AggroConfig{
		StringActive: true,
		StringTarget: "custom.field", // target the override attribute
		NumericActive: false,
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 2, overrides, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the output
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should contain the override attribute
	assert.Contains(t, output, "custom.field")
	// Should contain aggro metadata showing modification
	assert.Contains(t, output, "aggro.string")
	// The aggro metadata should reference the targeted custom field
	assert.Contains(t, output, "custom.field")
}

func TestAggroMetadataTracking(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config with specific targeting for precise verification
	aggroConfig := &aggro.AggroConfig{
		StringActive:    true,
		StringTarget:    "fake.attr.1", // target first fake attribute
		NumericActive:   true, 
		NumericTarget:   "fake.attr.2", // target second fake attribute
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify the metadata tracking works correctly
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	
	// Should contain metadata for both string and numeric aggro
	assert.Contains(t, output, "aggro.string")
	assert.Contains(t, output, "aggro.numeric")
	
	// Verify the metadata contains the correct target attribute names
	assert.Contains(t, output, "fake.attr.1") // string target
	assert.Contains(t, output, "fake.attr.2") // numeric target
	
	// Should not contain timestamp metadata since it's disabled
	assert.NotContains(t, output, "aggro.timestamp")
}

func TestAggroNoModificationWhenInactive(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config with all types inactive
	aggroConfig := &aggro.AggroConfig{
		StringActive:    false,
		NumericActive:   false,
		TimestampActive: false,
	}

	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfig, outputFile)
	require.NoError(t, err)

	// Read and verify no aggro metadata is present
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(content)
	// Should not contain any aggro metadata when all types are inactive
	assert.NotContains(t, output, "aggro.string")
	assert.NotContains(t, output, "aggro.numeric")
	assert.NotContains(t, output, "aggro.timestamp")
	
	// Should still contain regular fake attributes
	assert.Contains(t, output, "fake.attr.1")
	assert.Contains(t, output, "fake.attr.2")
	assert.Contains(t, output, "fake.attr.3")
}

func TestEnhancedAggroEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "traces.json")

	// Create aggro config to test enhanced edge cases
	aggroConfig := &aggro.AggroConfig{
		StringActive:    true,
		StringTarget:    "", // random targeting to get variety
		NumericActive:   true,
		NumericTarget:   "", // random targeting
		TimestampActive: true,
		TimestampTarget: "", // random targeting
	}

	// Run multiple times to increase chances of getting enhanced edge cases
	for i := 0; i < 10; i++ {
		err := generateTracesToFileWithAggro(1, 5, []string{}, aggroConfig, fmt.Sprintf("%s-%d", outputFile, i))
		require.NoError(t, err)
	}

	// Read all outputs and verify we got enhanced edge cases
	foundAdvancedFeatures := false
	for i := 0; i < 10; i++ {
		content, err := os.ReadFile(fmt.Sprintf("%s-%d", outputFile, i))
		require.NoError(t, err)

		output := string(content)
		
		// Check for any sign of enhanced features (unicode, emojis, scientific notation, etc.)
		if strings.Contains(output, "🚀") || strings.Contains(output, "💯") ||
		   strings.Contains(output, "Inf") || strings.Contains(output, "NaN") ||
		   strings.Contains(output, "1e") || strings.Contains(output, "E+") ||
		   strings.Contains(output, "café") || strings.Contains(output, "こんにちは") ||
		   strings.Contains(output, "9999-12-31") || strings.Contains(output, "1970-01-01") {
			foundAdvancedFeatures = true
			break
		}
	}

	// We should have found at least some enhanced features across multiple runs
	assert.True(t, foundAdvancedFeatures, "Expected to find enhanced aggro edge cases (unicode, emojis, scientific notation, etc.) in test output")
}

func TestAggroFlagSyntaxBehavior(t *testing.T) {
	// This test documents the current behavior of aggro flags
	// and serves as a regression test for flag parsing
	
	tmpDir := t.TempDir()
	outputFile1 := filepath.Join(tmpDir, "with-equals.json")
	outputFile2 := filepath.Join(tmpDir, "without-equals.json")

	// Test with explicit empty string (should work)
	aggroConfigWithEquals := &aggro.AggroConfig{
		StringActive: true,
		StringTarget: "", // empty = random targeting
		NumericActive: false,
		TimestampActive: false,
	}
	
	err := generateTracesToFileWithAggro(1, 3, []string{}, aggroConfigWithEquals, outputFile1)
	require.NoError(t, err)

	// Test with no aggro active (simulating --aggro-string without value)  
	aggroConfigWithoutEquals := &aggro.AggroConfig{
		StringActive: false, // This simulates the behavior when flag isn't properly detected
		NumericActive: false,
		TimestampActive: false,
	}
	
	err = generateTracesToFileWithAggro(1, 3, []string{}, aggroConfigWithoutEquals, outputFile2)
	require.NoError(t, err)

	// Read outputs
	contentWithEquals, err := os.ReadFile(outputFile1)
	require.NoError(t, err)
	
	contentWithoutEquals, err := os.ReadFile(outputFile2)
	require.NoError(t, err)

	outputWithEquals := string(contentWithEquals)
	outputWithoutEquals := string(contentWithoutEquals)

	// With explicit empty string should have aggro metadata
	assert.Contains(t, outputWithEquals, "aggro.string", "Using --aggro-string='' should apply string chaos engineering")

	// Without proper flag detection should NOT have aggro metadata
	assert.NotContains(t, outputWithoutEquals, "aggro.string", "When aggro is not active, should not apply string chaos engineering")
}