package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// getEnvFloat reads an environment variable and parses it as float64.
// If the value is missing or invalid, it returns the provided default.
func GetEnvFloat(key string, defaultVal float64) float64 {
	valStr := strings.TrimSpace(os.Getenv(key))
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Printf("Warning: invalid float for %s = %q, using default: %v", key, valStr, defaultVal)
		return defaultVal
	}

	return val
}
