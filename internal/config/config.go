package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialize sets up viper for configuration file and flag binding
func Initialize(rootCmd *cobra.Command) {
	// Allow viper to read from environment variables
	viper.AutomaticEnv()
	
	// Set the configuration file name and search paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.otel-datagen")
	viper.AddConfigPath("/etc/otel-datagen")
	
	// This will be called after flags are parsed
	cobra.OnInitialize(func() { initConfig(rootCmd) })
}

// initConfig reads in config file and ENV variables if set.
func initConfig(rootCmd *cobra.Command) {
	// Read config file if --config flag is provided
	if cfgFile := rootCmd.PersistentFlags().Lookup("config"); cfgFile != nil {
		if cfgFile.Value.String() != "" {
			viper.SetConfigFile(cfgFile.Value.String())
		}
	}
	
	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional, only log if explicitly provided
		if cfgFile := rootCmd.PersistentFlags().Lookup("config"); cfgFile != nil && cfgFile.Value.String() != "" {
			fmt.Printf("Warning: Could not read config file: %v\n", err)
		}
	}
	
	// Bind flags to viper for precedence handling
	bindFlags(rootCmd)
}

// bindFlags binds cobra flags to viper for configuration precedence
func bindFlags(rootCmd *cobra.Command) {
	// Find command references
	generateCmd, _, _ := rootCmd.Find([]string{"generate"})
	if generateCmd == nil {
		return
	}
	
	tracesCmd, _, _ := generateCmd.Find([]string{"traces"})
	logsCmd, _, _ := generateCmd.Find([]string{"logs"})
	metricsCmd, _, _ := generateCmd.Find([]string{"metrics"})
	
	// Global flags
	viper.BindPFlag("resource-attr", rootCmd.PersistentFlags().Lookup("resource-attr"))
	viper.BindPFlag("otlp-endpoint", rootCmd.PersistentFlags().Lookup("otlp-endpoint"))
	viper.BindPFlag("stdout", rootCmd.PersistentFlags().Lookup("stdout"))
	
	// Traces-specific flags  
	if tracesCmd != nil {
		viper.BindPFlag("generate.traces.num_traces", tracesCmd.Flags().Lookup("num-traces"))
		viper.BindPFlag("generate.traces.num_spans", tracesCmd.Flags().Lookup("num-spans"))
		viper.BindPFlag("generate.traces.num_attributes", tracesCmd.Flags().Lookup("num-attributes"))
		viper.BindPFlag("generate.traces.override_attr", tracesCmd.Flags().Lookup("override-attr"))
		viper.BindPFlag("generate.traces.aggro_timestamp", tracesCmd.Flags().Lookup("aggro-timestamp"))
		viper.BindPFlag("generate.traces.aggro_numeric", tracesCmd.Flags().Lookup("aggro-numeric"))
		viper.BindPFlag("generate.traces.aggro_string", tracesCmd.Flags().Lookup("aggro-string"))
	}
	
	// Logs-specific flags
	if logsCmd != nil {
		viper.BindPFlag("generate.logs.num_logs", logsCmd.Flags().Lookup("num-logs"))
		viper.BindPFlag("generate.logs.num_attributes", logsCmd.Flags().Lookup("num-attributes"))
		viper.BindPFlag("generate.logs.override_attr", logsCmd.Flags().Lookup("override-attr"))
		viper.BindPFlag("generate.logs.aggro_timestamp", logsCmd.Flags().Lookup("aggro-timestamp"))
		viper.BindPFlag("generate.logs.aggro_numeric", logsCmd.Flags().Lookup("aggro-numeric"))
		viper.BindPFlag("generate.logs.aggro_string", logsCmd.Flags().Lookup("aggro-string"))
	}
	
	// Metrics-specific flags
	if metricsCmd != nil {
		viper.BindPFlag("generate.metrics.num_metrics", metricsCmd.Flags().Lookup("num-metrics"))
		viper.BindPFlag("generate.metrics.metric_type", metricsCmd.Flags().Lookup("metric-type"))
		viper.BindPFlag("generate.metrics.metric_name", metricsCmd.Flags().Lookup("metric-name"))
		viper.BindPFlag("generate.metrics.counter_min", metricsCmd.Flags().Lookup("counter-min"))
		viper.BindPFlag("generate.metrics.counter_max", metricsCmd.Flags().Lookup("counter-max"))
		viper.BindPFlag("generate.metrics.aggro_timestamp", metricsCmd.Flags().Lookup("aggro-timestamp"))
		viper.BindPFlag("generate.metrics.aggro_numeric", metricsCmd.Flags().Lookup("aggro-numeric"))
		viper.BindPFlag("generate.metrics.aggro_string", metricsCmd.Flags().Lookup("aggro-string"))
	}
}