package envar

import "os"

// GetStringWithFallback returns the value of the environment variable with a fallback if not found.
func GetStringWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
