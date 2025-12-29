package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// ParseTimestamp converts a JWT timestamp (int64 or float64) to time.Time
func ParseTimestamp(claim interface{}) *time.Time {
	if claim == nil {
		return nil
	}

	var timestamp int64
	switch v := claim.(type) {
	case float64:
		timestamp = int64(v)
	case int64:
		timestamp = v
	case int:
		timestamp = int64(v)
	default:
		return nil
	}

	t := time.Unix(timestamp, 0)
	return &t
}

// IsTokenExpired checks if a token is expired based on exp claim
func IsTokenExpired(exp interface{}) bool {
	expTime := ParseTimestamp(exp)
	if expTime == nil {
		return false
	}
	return time.Now().After(*expTime)
}

// FormatTime formats a time.Time to a human-readable string
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05 MST")
}

// DecodeBase64URL decodes a base64url encoded string
func DecodeBase64URL(s string) ([]byte, error) {
	// Add padding if necessary
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(s)
}

// ParseJSONToClaims parses JSON data to a map
func ParseJSONToClaims(data []byte) (map[string]interface{}, error) {
	var claims map[string]interface{}
	err := json.Unmarshal(data, &claims)
	return claims, err
}

// GetStringClaim safely extracts a string claim
func GetStringClaim(claims map[string]interface{}, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
