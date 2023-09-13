package envar

import (
	"os"
	"strconv"
)

// GetIntWithFallback returns the value of the environment variable with a fallback if not found.
func GetIntWithFallback(value int, keys ...string) int {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			v, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}

			return v
		}
	}

	return value
}
