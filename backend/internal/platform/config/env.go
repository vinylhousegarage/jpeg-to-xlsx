package config

import (
	"fmt"
	"os"
	"strings"
)

func loadRequiredEnv(
	name string,
) (
	string,
	error,
) {
	value := os.Getenv(name)

	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf(
			"%s is required",
			name,
		)
	}

	return value, nil
}
