package envar

import (
	"os"
	"strconv"
)

// GetIntWithFallback returns the value of the environment variable with a fallback if not found.
func GetIntWithFallback(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}

		return v
	}

	return fallback
}
