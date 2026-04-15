package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot of the given env map to the specified file path.
func Save(path string, source string, env map[string]string) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Env:       env,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: failed to marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: failed to write file %q: %w", path, err)
	}

	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to read file %q: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: failed to parse file %q: %w", path, err)
	}

	return &snap, nil
}
