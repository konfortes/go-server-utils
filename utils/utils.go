package utils

import (
	"os"
)

// GetEnvOr gets an environment variable or returns ifNotFound value
func GetEnvOr(env, ifNotFound string) string {
	foundEnv, found := os.LookupEnv(env)

	if found {
		return foundEnv
	}

	return ifNotFound
}
