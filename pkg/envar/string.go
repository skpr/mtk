package envar

import "os"

// GetStringWithFallback returns the value of the environment variable with a fallback if not found.
func GetStringWithFallback(value string, keys ...string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
	}

	return value
}
