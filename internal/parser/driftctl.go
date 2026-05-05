package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// ResourceStatus describes whether a resource is changed, missing, or unmanaged.
type ResourceStatus string

const (
	StatusChanged   ResourceStatus = "changed"
	StatusMissing   ResourceStatus = "missing"
	StatusUnmanaged ResourceStatus = "unmanaged"
)

// Resource represents a single cloud resource identified in the drift report.
type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Difference captures a single drifted resource and its status.
type Difference struct {
	Res    Resource       `json:"res"`
	Status ResourceStatus `json:"status"`
}

// Summary holds high-level statistics from a drift scan.
type Summary struct {
	TotalResources   int `json:"total_resources"`
	DriftedResources int `json:"drifted_resources"`
	UndriftedCount   int `json:"undrifted_count"`
}

// DriftReport is the top-level structure of a driftctl JSON report.
type DriftReport struct {
	Summary     Summary      `json:"summary"`
	Differences []Difference `json:"differences"`
}

// ParseFile reads a driftctl JSON report from the given file path and returns
// a parsed DriftReport or an error.
func ParseFile(path string) (*DriftReport, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading report file %q: %w", path, err)
	}

	var report DriftReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parsing report JSON: %w", err)
	}

	return &report, nil
}
