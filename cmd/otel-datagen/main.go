package main

import (
	"fmt"
	"log"
	"os"

	"github.com/david/otel-datagen/internal/aggro"
	"github.com/david/otel-datagen/internal/config"
	"github.com/david/otel-datagen/internal/generators"
	"github.com/david/otel-datagen/internal/timestamps"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "otel-datagen",
	Short: "Generate synthetic OpenTelemetry data with intelligent randomness powered by the Antithesis SDK",
}

// parseTimestampConfig parses timestamp flags and returns a configuration
func parseTimestampConfig(cmd *cobra.Command) (*timestamps.TimestampConfig, error) {
	timestampStart, _ := cmd.Flags().GetString("timestamp-start")
	timestampSpacing, _ := cmd.Flags().GetString("timestamp-spacing")
	
	return timestamps.ParseTimestampConfig(timestampStart, timestampSpacing)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenTelemetry signals",
}

var tracesCmd = &cobra.Command{
	Use:   "traces",
	Short: "Generate trace data",
	Run: func(cmd *cobra.Command, args []string) {
		// Get values from viper (includes config file with CLI flag precedence)
		numTraces := viper.GetInt("generate.traces.num_traces")
		if numTraces == 0 {
			numTraces, _ = cmd.Flags().GetInt("num-traces")
		}
		
		numSpans := viper.GetInt("generate.traces.num_spans")
		if numSpans == 0 {
			numSpans, _ = cmd.Flags().GetInt("num-spans")
		}
		
		numAttributes := viper.GetInt("generate.traces.num_attributes")
		if numAttributes == 0 {
			numAttributes, _ = cmd.Flags().GetInt("num-attributes")
		}
		
		overrideAttrs := viper.GetStringSlice("generate.traces.override_attr")
		if len(overrideAttrs) == 0 {
			overrideAttrs, _ = cmd.Flags().GetStringSlice("override-attr")
		}
		
		
		// Parse new aggro configuration
		_ = aggro.ParseAggroConfig("traces")
		
		// Get resource attributes - CLI flags take precedence over config file
		resourceAttrs, _ := cmd.Root().PersistentFlags().GetStringSlice("resource-attr")
		
		// If no CLI resource attributes, get from config file format
		if len(resourceAttrs) == 0 {
			if resourceMap := viper.GetStringMapString("resource"); len(resourceMap) > 0 {
				for key, value := range resourceMap {
					resourceAttrs = append(resourceAttrs, key+"="+value)
				}
			}
		}
		
		otlpEndpoint := viper.GetString("otlp-endpoint")
		if otlpEndpoint == "" {
			otlpEndpoint, _ = cmd.Root().PersistentFlags().GetString("otlp-endpoint")
		}
		
		stdoutEnabled := viper.GetBool("stdout")
		if !stdoutEnabled {
			stdoutEnabled, _ = cmd.Root().PersistentFlags().GetBool("stdout")
		}
		
		// Automatically enable stdout when no OTLP endpoint is specified
		if otlpEndpoint == "" {
			stdoutEnabled = true
		}
		
		// Parse timestamp configuration
		timestampConfig, err := parseTimestampConfig(cmd.Parent())
		if err != nil {
			log.Fatalf("Error parsing timestamp configuration: %v", err)
		}
		
		generators.GenerateTraces(numTraces, numSpans, numAttributes, overrideAttrs, resourceAttrs, otlpEndpoint, stdoutEnabled, timestampConfig)
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Generate log data",
	Run: func(cmd *cobra.Command, args []string) {
		// Get values from viper (includes config file with CLI flag precedence)
		numLogs := viper.GetInt("generate.logs.num_logs")
		if numLogs == 0 {
			numLogs, _ = cmd.Flags().GetInt("num-logs")
		}
		
		numAttributes := viper.GetInt("generate.logs.num_attributes")
		if numAttributes == 0 {
			numAttributes, _ = cmd.Flags().GetInt("num-attributes")
		}
		
		overrideAttrs := viper.GetStringSlice("generate.logs.override_attr")
		if len(overrideAttrs) == 0 {
			overrideAttrs, _ = cmd.Flags().GetStringSlice("override-attr")
		}
		
		
		// Parse new aggro configuration
		_ = aggro.ParseAggroConfig("logs")
		
		// Get global flags - resource attributes and OTLP endpoint
		resourceAttrs, _ := cmd.Root().PersistentFlags().GetStringSlice("resource-attr")
		
		// If no CLI resource attributes, get from config file format
		if len(resourceAttrs) == 0 {
			if resourceMap := viper.GetStringMapString("resource"); len(resourceMap) > 0 {
				for key, value := range resourceMap {
					resourceAttrs = append(resourceAttrs, key+"="+value)
				}
			}
		}
		
		otlpEndpoint := viper.GetString("otlp-endpoint")
		if otlpEndpoint == "" {
			otlpEndpoint, _ = cmd.Root().PersistentFlags().GetString("otlp-endpoint")
		}
		
		stdoutEnabled := viper.GetBool("stdout")
		if !stdoutEnabled {
			stdoutEnabled, _ = cmd.Root().PersistentFlags().GetBool("stdout")
		}
		
		// Automatically enable stdout when no OTLP endpoint is specified
		if otlpEndpoint == "" {
			stdoutEnabled = true
		}
		
		// Default to gRPC protocol for now
		otlpProtocol := "grpc"
		
		// Parse timestamp configuration
		timestampConfig, err := parseTimestampConfig(cmd.Parent())
		if err != nil {
			log.Fatalf("Error parsing timestamp configuration: %v", err)
		}
		
		generators.GenerateLogs(numLogs, numAttributes, overrideAttrs, resourceAttrs, otlpEndpoint, otlpProtocol, stdoutEnabled, timestampConfig)
	},
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Generate metric data",
	Run: func(cmd *cobra.Command, args []string) {
		// Get values from viper (includes config file with CLI flag precedence)
		numMetrics := viper.GetInt("generate.metrics.num_metrics")
		if numMetrics == 0 {
			numMetrics, _ = cmd.Flags().GetInt("num-metrics")
		}
		
		metricType := viper.GetString("generate.metrics.metric_type")
		if metricType == "" {
			metricType, _ = cmd.Flags().GetString("metric-type")
		}
		
		metricName := viper.GetString("generate.metrics.metric_name")
		if metricName == "" {
			metricName, _ = cmd.Flags().GetString("metric-name")
		}
		
		counterMin := viper.GetInt("generate.metrics.counter_min")
		if counterMin == 0 {
			counterMin, _ = cmd.Flags().GetInt("counter-min")
		}
		
		counterMax := viper.GetInt("generate.metrics.counter_max")
		if counterMax == 0 {
			counterMax, _ = cmd.Flags().GetInt("counter-max")
		}
		
		
		// Parse new aggro configuration
		_ = aggro.ParseAggroConfig("metrics")
		
		// Get global flags - resource attributes and OTLP endpoint
		resourceAttrs, _ := cmd.Root().PersistentFlags().GetStringSlice("resource-attr")
		
		// If no CLI resource attributes, get from config file format
		if len(resourceAttrs) == 0 {
			if resourceMap := viper.GetStringMapString("resource"); len(resourceMap) > 0 {
				for key, value := range resourceMap {
					resourceAttrs = append(resourceAttrs, key+"="+value)
				}
			}
		}
		
		otlpEndpoint := viper.GetString("otlp-endpoint")
		if otlpEndpoint == "" {
			otlpEndpoint, _ = cmd.Root().PersistentFlags().GetString("otlp-endpoint")
		}
		
		stdoutEnabled := viper.GetBool("stdout")
		if !stdoutEnabled {
			stdoutEnabled, _ = cmd.Root().PersistentFlags().GetBool("stdout")
		}
		
		// Automatically enable stdout when no OTLP endpoint is specified
		if otlpEndpoint == "" {
			stdoutEnabled = true
		}
		
		// Default to gRPC protocol for now
		otlpProtocol := "grpc"
		
		// Parse timestamp configuration
		timestampConfig, err := parseTimestampConfig(cmd.Parent())
		if err != nil {
			log.Fatalf("Error parsing timestamp configuration: %v", err)
		}
		
		generators.GenerateMetrics(numMetrics, metricType, metricName, counterMin, counterMax, resourceAttrs, otlpEndpoint, otlpProtocol, stdoutEnabled, timestampConfig)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(tracesCmd)
	generateCmd.AddCommand(logsCmd)
	generateCmd.AddCommand(metricsCmd)
	
	// Set up flags using config package
	config.SetupFlags(rootCmd, generateCmd, tracesCmd, logsCmd, metricsCmd)
	
	// Set up viper configuration
	config.Initialize(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}