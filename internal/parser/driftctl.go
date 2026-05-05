package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// DriftResource represents a single drifted resource from driftctl output.
type DriftResource struct {
	ResourceID   string            `json:"id"`
	ResourceType string            `json:"type"`
	Source       string            `json:"source"`
	Differences  []Difference      `json:"differences,omitempty"`
	Attrs        map[string]string `json:"attrs,omitempty"`
}

// Difference captures a single attribute-level drift.
type Difference struct {
	Attribute string `json:"attribute"`
	Wanted    string `json:"wanted"`
	Found     string `json:"found"`
}

// DriftReport is the top-level structure parsed from a driftctl JSON output file.
type DriftReport struct {
	Summary struct {
		TotalResources  int `json:"total_resources"`
		DriftedCount    int `json:"total_drifted"`
		UnmanagedCount  int `json:"total_unmanaged"`
		DeletedCount    int `json:"total_deleted"`
	} `json:"summary"`
	Drifted   []DriftResource `json:"differences"`
	Unmanaged []DriftResource `json:"unmanaged"`
	Deleted   []DriftResource `json:"deleted"`
}

// ParseFile reads and parses a driftctl JSON report from the given file path.
func ParseFile(path string) (*DriftReport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: opening file %q: %w", path, err)
	}
	defer f.Close()

	var report DriftReport
	if err := json.NewDecoder(f).Decode(&report); err != nil {
		return nil, fmt.Errorf("parser: decoding JSON from %q: %w", path, err)
	}
	return &report, nil
}
