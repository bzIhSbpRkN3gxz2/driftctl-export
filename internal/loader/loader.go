// Package loader provides utilities for loading driftctl report files
// from the filesystem, supporting both single file and glob patterns.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/driftctl-export/internal/parser"
	"github.com/cloudskiff/driftctl/pkg/resource"
)

// Result holds the outcome of loading one or more report files.
type Result struct {
	// Reports maps each resolved file path to its parsed report.
	Reports map[string]*resource.DriftctlScanReport
	// Errors maps file paths to any errors encountered during loading.
	Errors map[string]error
}

// HasErrors returns true if any files failed to load or parse.
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// LoadFile loads and parses a single driftctl JSON report from the given path.
func LoadFile(path string) (*resource.DriftctlScanReport, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("report file not found: %s", path)
	}
	return parser.ParseFile(path)
}

// LoadGlob loads all driftctl JSON reports matching the given glob pattern.
// Files that fail to parse are recorded in Result.Errors rather than
// causing the entire load to fail.
func LoadGlob(pattern string) (*Result, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no files matched pattern: %s", pattern)
	}

	result := &Result{
		Reports: make(map[string]*resource.DriftctlScanReport),
		Errors:  make(map[string]error),
	}

	for _, match := range matches {
		report, parseErr := parser.ParseFile(match)
		if parseErr != nil {
			result.Errors[match] = parseErr
			continue
		}
		result.Reports[match] = report
	}

	return result, nil
}
