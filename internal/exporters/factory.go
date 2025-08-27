package exporters

import (
	"context"
	"io"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// ExporterConfig holds configuration for exporters
type ExporterConfig struct {
	OTLPEndpoint string
	Protocol     string // "grpc" or "http"
	Insecure     bool
	Headers      map[string]string
	StdoutEnabled bool
}

// CreateDualTraceExporters creates both OTLP and console trace exporters when OTLP endpoint is specified
func CreateDualTraceExporters(ctx context.Context, config ExporterConfig, writer io.Writer) ([]trace.SpanExporter, error) {
	var exporters []trace.SpanExporter
	
	if config.OTLPEndpoint != "" {
		// Create console exporter only if stdout is enabled
		if config.StdoutEnabled {
			consoleOpts := []stdouttrace.Option{
				stdouttrace.WithPrettyPrint(),
			}
			if writer != nil {
				consoleOpts = append(consoleOpts, stdouttrace.WithWriter(writer))
			}
			consoleExporter, err := stdouttrace.New(consoleOpts...)
			if err != nil {
				return nil, err
			}
			exporters = append(exporters, consoleExporter)
		}
		
		// Create OTLP exporter
		otlpOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(config.OTLPEndpoint),
		}
		
		if config.Insecure {
			otlpOpts = append(otlpOpts, otlptracegrpc.WithInsecure())
		}
		
		if len(config.Headers) > 0 {
			otlpOpts = append(otlpOpts, otlptracegrpc.WithHeaders(config.Headers))
		}
		
		otlpExporter, err := otlptracegrpc.New(ctx, otlpOpts...)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, otlpExporter)
	} else {
		// Single console exporter when no OTLP endpoint
		opts := []stdouttrace.Option{
			stdouttrace.WithPrettyPrint(),
		}
		if writer != nil {
			opts = append(opts, stdouttrace.WithWriter(writer))
		}
		
		exporter, err := stdouttrace.New(opts...)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, exporter)
	}
	
	return exporters, nil
}

// CreateDualLogExporters creates both OTLP and console log exporters when OTLP endpoint is specified
func CreateDualLogExporters(ctx context.Context, config ExporterConfig, writer io.Writer) ([]sdklog.Exporter, error) {
	var exporters []sdklog.Exporter
	
	if config.OTLPEndpoint != "" {
		// Create console exporter only if stdout is enabled
		if config.StdoutEnabled {
			consoleOpts := []stdoutlog.Option{
				stdoutlog.WithPrettyPrint(),
			}
			if writer != nil {
				consoleOpts = append(consoleOpts, stdoutlog.WithWriter(writer))
			}
			consoleExporter, err := stdoutlog.New(consoleOpts...)
			if err != nil {
				return nil, err
			}
			exporters = append(exporters, consoleExporter)
		}
		
		// Create OTLP exporter
		var otlpExporter sdklog.Exporter
		var err error
		
		if config.Protocol == "http" {
			opts := []otlploghttp.Option{
				otlploghttp.WithEndpoint(config.OTLPEndpoint),
			}
			
			if config.Insecure {
				opts = append(opts, otlploghttp.WithInsecure())
			}
			
			if len(config.Headers) > 0 {
				opts = append(opts, otlploghttp.WithHeaders(config.Headers))
			}
			
			otlpExporter, err = otlploghttp.New(ctx, opts...)
		} else {
			// Default to gRPC
			opts := []otlploggrpc.Option{
				otlploggrpc.WithEndpoint(config.OTLPEndpoint),
			}
			
			if config.Insecure {
				opts = append(opts, otlploggrpc.WithInsecure())
			}
			
			if len(config.Headers) > 0 {
				opts = append(opts, otlploggrpc.WithHeaders(config.Headers))
			}
			
			otlpExporter, err = otlploggrpc.New(ctx, opts...)
		}
		
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, otlpExporter)
	} else {
		// Single console exporter when no OTLP endpoint
		opts := []stdoutlog.Option{
			stdoutlog.WithPrettyPrint(),
		}
		if writer != nil {
			opts = append(opts, stdoutlog.WithWriter(writer))
		}
		
		exporter, err := stdoutlog.New(opts...)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, exporter)
	}
	
	return exporters, nil
}

// CreateDualMetricExporters creates both OTLP and console metric exporters when OTLP endpoint is specified
func CreateDualMetricExporters(ctx context.Context, config ExporterConfig, writer io.Writer) ([]metric.Exporter, error) {
	var exporters []metric.Exporter
	
	if config.OTLPEndpoint != "" {
		// Create console exporter only if stdout is enabled
		if config.StdoutEnabled {
			consoleOpts := []stdoutmetric.Option{
				stdoutmetric.WithPrettyPrint(),
			}
			if writer != nil {
				consoleOpts = append(consoleOpts, stdoutmetric.WithWriter(writer))
			}
			consoleExporter, err := stdoutmetric.New(consoleOpts...)
			if err != nil {
				return nil, err
			}
			exporters = append(exporters, consoleExporter)
		}
		
		// Create OTLP exporter
		var otlpExporter metric.Exporter
		var err error
		
		if config.Protocol == "http" {
			opts := []otlpmetrichttp.Option{
				otlpmetrichttp.WithEndpoint(config.OTLPEndpoint),
			}
			
			if config.Insecure {
				opts = append(opts, otlpmetrichttp.WithInsecure())
			}
			
			if len(config.Headers) > 0 {
				opts = append(opts, otlpmetrichttp.WithHeaders(config.Headers))
			}
			
			otlpExporter, err = otlpmetrichttp.New(ctx, opts...)
		} else {
			// Default to gRPC
			opts := []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithEndpoint(config.OTLPEndpoint),
			}
			
			if config.Insecure {
				opts = append(opts, otlpmetricgrpc.WithInsecure())
			}
			
			if len(config.Headers) > 0 {
				opts = append(opts, otlpmetricgrpc.WithHeaders(config.Headers))
			}
			
			otlpExporter, err = otlpmetricgrpc.New(ctx, opts...)
		}
		
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, otlpExporter)
	} else {
		// Single console exporter when no OTLP endpoint
		opts := []stdoutmetric.Option{
			stdoutmetric.WithPrettyPrint(),
		}
		if writer != nil {
			opts = append(opts, stdoutmetric.WithWriter(writer))
		}
		
		exporter, err := stdoutmetric.New(opts...)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, exporter)
	}
	
	return exporters, nil
}