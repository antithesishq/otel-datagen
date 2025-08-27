package config

import (
	"github.com/david/otel-datagen/internal/randomness"
	"github.com/spf13/cobra"
)

// SetupFlags adds all CLI flags to the commands
func SetupFlags(rootCmd, generateCmd, tracesCmd, logsCmd, metricsCmd *cobra.Command) {
	// Global flags for all commands
	rootCmd.PersistentFlags().String("config", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringSlice("resource-attr", []string{}, "Set resource attributes (key=value)")
	rootCmd.PersistentFlags().String("otlp-endpoint", "", "OTLP endpoint URL (if set, exports to OTLP)")
	rootCmd.PersistentFlags().Bool("stdout", false, "Output to stdout console (automatically enabled when no OTLP endpoint is set)")

	// Timestamp control flags
	generateCmd.PersistentFlags().String("timestamp-start", "", "Start timestamp for data generation (ISO 8601 or relative like '-5m', '-1h')")
	generateCmd.PersistentFlags().String("timestamp-spacing", "0s", "Duration between consecutive data points (e.g., '30s', '1m')")

	// Traces-specific flags
	tracesCmd.Flags().Int("num-traces", 1, "Number of traces to generate")
	tracesCmd.Flags().Int("num-spans", randomness.Intn(5)+1, "Number of spans to generate per trace")
	tracesCmd.Flags().Int("num-attributes", 5, "Number of additional random attributes to add")
	tracesCmd.Flags().StringSlice("override-attr", []string{}, "Override specific attributes (key=value)")
	tracesCmd.Flags().String("aggro-timestamp", "", "Apply timestamp chaos engineering (empty=random, 'attr'=target specific attribute)")
	tracesCmd.Flags().String("aggro-numeric", "", "Apply numeric chaos engineering (empty=random, 'attr'=target specific attribute)")  
	tracesCmd.Flags().String("aggro-string", "", "Apply string chaos engineering (empty=random, 'attr'=target specific attribute)")

	// Logs-specific flags
	logsCmd.Flags().Int("num-logs", 1, "Number of log records to generate")
	logsCmd.Flags().Int("num-attributes", 5, "Number of additional random attributes to add")
	logsCmd.Flags().StringSlice("override-attr", []string{}, "Override specific attributes (key=value)")
	logsCmd.Flags().String("aggro-timestamp", "", "Apply timestamp chaos engineering (empty=random, 'attr'=target specific attribute)")
	logsCmd.Flags().String("aggro-numeric", "", "Apply numeric chaos engineering (empty=random, 'attr'=target specific attribute)")  
	logsCmd.Flags().String("aggro-string", "", "Apply string chaos engineering (empty=random, 'attr'=target specific attribute)")

	// Metrics-specific flags
	metricsCmd.Flags().Int("num-metrics", 5, "Number of metric data points to generate")
	metricsCmd.Flags().String("metric-type", "counter", "Type of metric: counter, float64-counter, histogram, int64-histogram, updowncounter, float64-updowncounter, gauge, float64-gauge, observable-counter, float64-observable-counter, observable-updowncounter, float64-observable-updowncounter, observable-gauge, float64-observable-gauge")
	metricsCmd.Flags().String("metric-name", "example_metric", "Name of the metric")
	metricsCmd.Flags().Int("counter-min", 1, "Minimum value for metrics")
	metricsCmd.Flags().Int("counter-max", 100, "Maximum value for metrics")
	metricsCmd.Flags().String("aggro-timestamp", "", "Apply timestamp chaos engineering (empty=random, 'attr'=target specific attribute)")
	metricsCmd.Flags().String("aggro-numeric", "", "Apply numeric chaos engineering (empty=random, 'attr'=target specific attribute)")  
	metricsCmd.Flags().String("aggro-string", "", "Apply string chaos engineering (empty=random, 'attr'=target specific attribute)")
}
