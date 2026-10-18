package identity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func LoadOrCreate(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		id := strings.TrimSpace(string(data))
		if id == "" {
			return "", fmt.Errorf("identity file is empty: %s", path)
		}
		if _, err := uuid.Parse(id); err != nil {
			return "", fmt.Errorf("invalid identity in %s: %w", path, err)
		}
		return id, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("read identity file: %w", err)
	}
	id := uuid.NewString()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", fmt.Errorf("create identity directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(id), 0700); err != nil {
		return "", fmt.Errorf("write identity file: %w", err)
	}
	return id, nil
}
