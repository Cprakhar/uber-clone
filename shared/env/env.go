package env

import (
	"os"
	"strconv"
	"time"
)

// GetString retrieves the value of the environment variable named by the key.
// If the variable is empty or not present, it returns the specified default value.
func GetString(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// GetInt retrieves the value of the environment variable named by the key and converts it to an integer.
// If the variable is empty, not present, or cannot be converted to an integer, it returns the specified default value.
func GetInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetBool retrieves the value of the environment variable named by the key and converts it to a boolean.
// If the variable is empty, not present, or cannot be converted to a boolean, it returns the specified default value.
func GetBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// GetDuration retrieves the value of the environment variable named by the key and converts it to a time.Duration.
// If the variable is empty, not present, or cannot be converted to a duration, it returns the specified default value.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	durationVal, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return durationVal
}
