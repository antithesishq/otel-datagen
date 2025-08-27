# otel-datagen

OpenTelemetry data generation CLI tool for testing observability pipelines with intelligent randomness powered by the Antithesis SDK.

## Getting Started

Build the binary:
```bash
go build -o otel-datagen ./cmd/otel-datagen
```

## Usage

Generate traces and output to console:
```bash
./otel-datagen generate traces
```

Generate multiple traces with multiple spans per trace:
```bash
./otel-datagen generate traces --num-traces=3 --num-spans=5
```

Generate traces with custom attributes:
```bash
./otel-datagen generate traces --num-attributes=7
```

Override specific attributes:
```bash
./otel-datagen generate traces --override-attr=service.name=my-test-service
```

Generate traces with chaos engineering (aggro testing):
```bash
# Apply random string chaos engineering
./otel-datagen generate traces --aggro-string=""

# Apply numeric chaos engineering to specific attribute
./otel-datagen generate traces --aggro-numeric="custom.metric"

# Apply timestamp chaos engineering
./otel-datagen generate traces --aggro-timestamp=""

# Apply multiple aggro types simultaneously
./otel-datagen generate traces --aggro-string="" --aggro-numeric="" --aggro-timestamp=""
```

Set custom resource attributes:
```bash
./otel-datagen generate traces --resource-attr service.name=my-service --resource-attr service.version=1.0.0
```

Export traces to OTLP endpoint:
```bash
./otel-datagen generate traces --otlp-endpoint http://localhost:4317
```

## Log Generation

Generate log records with realistic data:

```bash
./otel-datagen generate logs
```

Generate multiple log records:
```bash
./otel-datagen generate logs --num-logs=10
```

Generate logs with custom attributes:
```bash
./otel-datagen generate logs --num-attributes=7
```

Generate logs with chaos engineering (aggro testing):
```bash
# Apply string chaos engineering to log messages
./otel-datagen generate logs --aggro-string="message"

# Apply random numeric chaos engineering
./otel-datagen generate logs --aggro-numeric=""

# Apply timestamp chaos engineering to specific attribute
./otel-datagen generate logs --aggro-timestamp="event.timestamp"
```

Override specific attributes:
```bash
./otel-datagen generate logs --override-attr=environment=production --override-attr=level=error
```

Set custom resource attributes for logs:
```bash
./otel-datagen generate logs --resource-attr service.name=my-log-service --resource-attr service.version=2.0.0
```

**Note**: OTLP log export is not yet available in the current OpenTelemetry Go SDK. When using `--otlp-endpoint` with logs, the tool will emit a warning and fall back to stdout output.

## Metrics Generation

Generate metrics with configurable types and values:

```bash
./otel-datagen generate metrics
```

Generate counter metrics:
```bash
./otel-datagen generate metrics --metric-type=counter --metric-name=request_count
```

Generate updowncounter metrics:
```bash
./otel-datagen generate metrics --metric-type=updowncounter --metric-name=queue_size
```

Generate gauge metrics:
```bash
./otel-datagen generate metrics --metric-type=gauge --metric-name=cpu_usage
```

Generate histogram metrics:
```bash
./otel-datagen generate metrics --metric-type=histogram --metric-name=response_time
```

Generate observable metrics:
```bash
./otel-datagen generate metrics --metric-type=observablecounter --metric-name=total_requests
```

Generate multiple metric data points:
```bash
./otel-datagen generate metrics --num-metrics=10
```

Configure counter value ranges:
```bash
./otel-datagen generate metrics --counter-min=1 --counter-max=1000
```

Generate metrics with chaos engineering (aggro testing):
```bash
# Apply numeric chaos engineering to metrics
./otel-datagen generate metrics --aggro-numeric=""

# Apply string chaos engineering to specific attribute
./otel-datagen generate metrics --aggro-string="service.instance"

# Apply timestamp chaos engineering
./otel-datagen generate metrics --aggro-timestamp=""
```

Set custom resource attributes for metrics:
```bash
./otel-datagen generate metrics --resource-attr service.name=metrics-service --resource-attr environment=prod
```

### Supported Metric Types

- **counter**: Int64 counter metrics with configurable min/max values (monotonic, always increasing)
- **float64counter**: Float64 counter metrics with configurable min/max values (monotonic, always increasing)
- **histogram**: Float64 histogram metrics for measuring distributions
- **int64histogram**: Int64 histogram metrics for measuring distributions
- **updowncounter**: Int64 updowncounter metrics that can increase or decrease
- **float64updowncounter**: Float64 updowncounter metrics that can increase or decrease
- **gauge**: Int64 gauge metrics for current values that can go up or down
- **float64gauge**: Float64 gauge metrics for current values that can go up or down
- **observablecounter**: Observable Int64 counter metrics (async callback-based)
- **observableupdowncounter**: Observable Int64 updowncounter metrics (async callback-based)
- **observablegauge**: Observable Int64 gauge metrics (async callback-based)
- **observablefloat64counter**: Observable Float64 counter metrics (async callback-based)
- **observablefloat64updowncounter**: Observable Float64 updowncounter metrics (async callback-based)
- **observablefloat64gauge**: Observable Float64 gauge metrics (async callback-based)

## Chaos Engineering (Aggro Testing)

The tool supports chaos engineering through "aggro" flags that inject edge case values into generated OpenTelemetry data to test system resilience:

### Aggro Flag Types

- **`--aggro-string[=attribute]`**: Injects naughty strings (from the Big List of Naughty Strings) 
- **`--aggro-numeric[=attribute]`**: Injects numeric aggro values (zero, negative, max values, etc.)
- **`--aggro-timestamp[=attribute]`**: Injects timestamp edge cases (epoch, far future, far past, etc.)

### Targeting Modes

Each aggro flag supports two modes:

1. **Random targeting** (empty string): `--aggro-string=""`
   - Randomly selects an attribute of the appropriate type to modify
   - Use when you want to test general system resilience

2. **Specific targeting**: `--aggro-string="user.name"`  
   - Targets a specific attribute name for modification
   - Use when testing specific attribute handling

### Examples

```bash
# Random string chaos engineering
./otel-datagen generate traces --aggro-string=""

# Target specific log message
./otel-datagen generate logs --aggro-string="message"

# Multiple aggro types simultaneously  
./otel-datagen generate traces --aggro-string="" --aggro-numeric="count" --aggro-timestamp=""
```

## Configuration File Support

You can configure all settings using a YAML configuration file with the `--config` flag:

```bash
./otel-datagen --config config.yaml generate traces
```

### Configuration File Structure

Example `config.yaml`:
```yaml
# Global resource attributes
resource:
  service.name: "my-service"
  service.version: "1.2.3"
  environment: "production"

# OTLP endpoint for remote export (optional)
otlp-endpoint: "http://localhost:4317"

# Generation settings
generate:
  traces:
    num_spans: 10
    num_attributes: 5
    aggro_string: ""              # Apply random string chaos engineering
    aggro_numeric: "custom.attr"  # Apply numeric chaos engineering to specific attribute
    aggro_timestamp: ""           # Apply random timestamp chaos engineering
    override_attr:
      - "custom.key=custom-value"
      - "another.key=another-value"
  logs:
    num_logs: 5
    num_attributes: 3
    aggro_string: "message"       # Apply string chaos engineering to log messages
    aggro_numeric: ""             # Apply random numeric chaos engineering
    override_attr:
      - "log.level=warn"
      - "environment=staging"
  metrics:
    num_metrics: 8
    metric_type: "histogram"
    metric_name: "custom_histogram"
    counter_min: 10
    counter_max: 500
    aggro_numeric: ""             # Apply numeric chaos engineering to metrics
```

### Configuration Precedence

Configuration values are resolved in the following order (highest to lowest precedence):
1. CLI flags (e.g., `--num-traces=2 --num-spans=5`)
2. Configuration file values
3. Default values

For example, if your config file sets `num_spans: 10` but you run `./otel-datagen --config config.yaml generate traces --num-spans=5`, the CLI flag value of `5` will be used. Similarly, you can set both `num_traces: 2` and `num_spans: 3` in your config file.

## Antithesis Integration

This tool uses the [Antithesis Go SDK](https://antithesis.com/docs/using_antithesis/sdk/go/) for intelligent randomness generation, which provides several benefits:

### Enhanced Testing with Antithesis Platform

When running on the Antithesis platform, the tool benefits from:
- **Intelligent edge case exploration**: Antithesis guides randomness toward interesting scenarios and aggro conditions
- **Reproducible test runs**: Deterministic randomness for debugging and issue reproduction  
- **Systematic coverage**: More thorough exploration of configuration combinations and value ranges
- **Enhanced aggro testing**: Better utilization of the "Big List of Naughty Strings" for resilience testing

### Seamless Fallback Behavior

When running outside the Antithesis environment (normal usage), the tool:
- Automatically falls back to `crypto/rand` for secure randomness
- Maintains identical functionality and performance
- Requires no additional configuration or setup
- Works exactly as it did before integration

### What Uses Antithesis Randomness

The following aspects of data generation now use Antithesis-guided randomness:
- **Aggro chaos engineering**: When and which attributes to target for chaos engineering
- **Aggro value selection**: Which specific naughty strings, numeric edge cases, or timestamp boundaries to inject
- **Metric value generation**: Counter, gauge, and histogram values within specified ranges
- **UpDownCounter negative values**: 30% probability of generating negative values  
- **Attribute generation**: Selection of fake data values for realistic attributes

### Running with Antithesis

To run this tool on the Antithesis platform, follow the standard Antithesis setup procedures. No code changes are required - the same binary works in both environments.

Example output:
```json
{
  "Name": "example-span",
  "SpanContext": {
    "TraceID": "1327a029bff2e83104e3e944cfa52a83",
    "SpanID": "1290a1595b86d3ea",
    "TraceFlags": "01"
  },
  "Resource": [
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "otel-datagen"
      }
    }
  ]
}
```