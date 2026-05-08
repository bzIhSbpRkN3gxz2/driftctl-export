// Package enricher provides functionality to annotate drift report resources
// with additional metadata such as tags, environment labels, and severity hints.
package enricher

import (
	"strings"

	"github.com/snyk/driftctl-export/internal/parser"
)

// Options controls which enrichment steps are applied.
type Options struct {
	// Environment label to attach to all resources (e.g. "production", "staging").
	Environment string

	// SeverityMap maps resource types to a severity string (e.g. "high", "low").
	SeverityMap map[string]string

	// DefaultSeverity is used when a resource type has no entry in SeverityMap.
	DefaultSeverity string
}

// Apply enriches the report in-place according to the provided Options.
// It returns the same report pointer for chaining convenience.
func Apply(report *parser.Report, opts Options) *parser.Report {
	if report == nil {
		return nil
	}

	env := strings.TrimSpace(opts.Environment)
	defSev := opts.DefaultSeverity
	if defSev == "" {
		defSev = "medium"
	}

	for i := range report.Summary.TotalUnmanagedResources {
		r := &report.UnmanagedResources[i]
		if env != "" {
			r.Environment = env
		}
		r.Severity = resolveSeverity(r.Type, opts.SeverityMap, defSev)
	}

	for i := range report.ManagedResources {
		r := &report.ManagedResources[i]
		if env != "" {
			r.Environment = env
		}
		r.Severity = resolveSeverity(r.Type, opts.SeverityMap, defSev)
	}

	return report
}

func resolveSeverity(resourceType string, severityMap map[string]string, defaultSev string) string {
	if severityMap == nil {
		return defaultSev
	}
	if sev, ok := severityMap[resourceType]; ok {
		return sev
	}
	return defaultSev
}
