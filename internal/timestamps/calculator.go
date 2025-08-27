package timestamps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimestampConfig holds configuration for timestamp generation
type TimestampConfig struct {
	StartTime time.Time
	Spacing   time.Duration
}

// ParseTimestampConfig parses timestamp-start and timestamp-spacing flags
func ParseTimestampConfig(timestampStart, timestampSpacing string) (*TimestampConfig, error) {
	// Parse spacing duration
	spacing, err := time.ParseDuration(timestampSpacing)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp-spacing '%s': %w", timestampSpacing, err)
	}

	// Parse start time
	var startTime time.Time
	if timestampStart == "" {
		// Default to current time
		startTime = time.Now()
	} else if strings.HasPrefix(timestampStart, "-") {
		// Relative time (e.g., "-5m", "-1h")
		relativeDuration, err := time.ParseDuration(timestampStart[1:]) // Remove the '-' prefix
		if err != nil {
			return nil, fmt.Errorf("invalid relative timestamp-start '%s': %w", timestampStart, err)
		}
		startTime = time.Now().Add(-relativeDuration)
	} else {
		// Try to parse as ISO 8601 timestamp
		startTime, err = parseTimestamp(timestampStart)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp-start '%s': %w", timestampStart, err)
		}
	}

	return &TimestampConfig{
		StartTime: startTime,
		Spacing:   spacing,
	}, nil
}

// CalculateTimestamp calculates the timestamp for the i-th data point
func (tc *TimestampConfig) CalculateTimestamp(index int) time.Time {
	return tc.StartTime.Add(time.Duration(index) * tc.Spacing)
}

// parseTimestamp attempts to parse various timestamp formats
func parseTimestamp(timestampStr string) (time.Time, error) {
	// List of supported timestamp formats
	formats := []string{
		time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	// Try parsing as Unix timestamp (seconds since epoch)
	if unixSeconds, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
		return time.Unix(unixSeconds, 0), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp format")
}