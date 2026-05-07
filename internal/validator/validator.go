// Package validator provides validation logic for driftctl report structures
// before they are processed by the export pipeline.
package validator

import (
	"errors"
	"fmt"

	"github.com/snyk/driftctl-export/internal/parser"
)

// ValidationError holds a list of issues found during report validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d issue(s): %v", len(e.Issues), e.Issues)
}

// Validate checks a parsed DriftReport for structural and semantic correctness.
// It returns a *ValidationError if any issues are found, or nil on success.
func Validate(report *parser.DriftReport) error {
	if report == nil {
		return errors.New("report must not be nil")
	}

	var issues []string

	for i, r := range report.Differences {
		if r.Res.ResourceType == "" {
			issues = append(issues, fmt.Sprintf("difference[%d]: missing resource type", i))
		}
		if r.Res.ResourceID == "" {
			issues = append(issues, fmt.Sprintf("difference[%d]: missing resource id", i))
		}
	}

	for i, r := range report.UnmanagedResources {
		if r.ResourceType == "" {
			issues = append(issues, fmt.Sprintf("unmanaged[%d]: missing resource type", i))
		}
		if r.ResourceID == "" {
			issues = append(issues, fmt.Sprintf("unmanaged[%d]: missing resource id", i))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
