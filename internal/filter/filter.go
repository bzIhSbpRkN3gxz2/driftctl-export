// Package filter provides utilities for filtering drift report resources
// based on configurable criteria such as resource type, drift kind, or namespace.
package filter

import "github.com/snyk/driftctl/pkg/resource"

// Options holds the filtering criteria applied to a drift report.
type Options struct {
	// ResourceTypes restricts output to the given Terraform resource types.
	// An empty slice means no restriction.
	ResourceTypes []string

	// OnlyManaged, when true, includes only resources that are managed by
	// Terraform (i.e. present in state).
	OnlyManaged bool

	// OnlyDrifted, when true, includes only resources that have at least one
	// changed attribute.
	OnlyDrifted bool
}

// ResourceSummary is a flattened view of a single resource used throughout
// the export pipeline.
type ResourceSummary struct {
	ID           string
	Type         string
	Source       string // "managed", "unmanaged", or "deleted"
	ChangedAttrs []string
}

// Apply filters the resources contained in report according to opts and
// returns a slice of ResourceSummary values ready for export.
func Apply(report *resource.ScanReport, opts Options) []ResourceSummary {
	var results []ResourceSummary

	append := func(id, resType, source string, changed []string) {
		if len(opts.ResourceTypes) > 0 && !containsString(opts.ResourceTypes, resType) {
			return
		}
		if opts.OnlyDrifted && len(changed) == 0 {
			return
		}
		results = append(results, ResourceSummary{
			ID:           id,
			Type:         resType,
			Source:       source,
			ChangedAttrs: changed,
		})
	}

	for _, r := range report.Managed {
		changed := changedAttributes(r)
		append(r.ResourceId(), r.ResourceType(), "managed", changed)
	}

	if !opts.OnlyManaged {
		for _, r := range report.Unmanaged {
			append(r.ResourceId(), r.ResourceType(), "unmanaged", nil)
		}
		for _, r := range report.Deleted {
			append(r.ResourceId(), r.ResourceType(), "deleted", nil)
		}
	}

	return results
}

func changedAttributes(r resource.Resource) []string {
	if differ, ok := r.(interface {
		Diff() map[string]interface{}
	}); ok {
		keys := make([]string, 0, len(differ.Diff()))
		for k := range differ.Diff() {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
