package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvVarAPIKey = "HACKSPLAINING_API_KEY"
	ConfigFile   = ".hacksplaining"
)

func LoadAPIKey() (string, error) {
	if key := os.Getenv(EnvVarAPIKey); key != "" {
		return key, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determining home directory: %w", err)
	}

	paths := []string{
		filepath.Join(home, ConfigFile),
		filepath.Join("/etc", ConfigFile),
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		key := strings.TrimSpace(string(data))
		if key != "" {
			return key, nil
		}
	}

	return "", fmt.Errorf(
		"API key not found. Set %s or create ~/%s or /etc/%s containing your key",
		EnvVarAPIKey, ConfigFile, ConfigFile,
	)
}
