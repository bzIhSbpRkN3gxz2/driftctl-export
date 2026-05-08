// Package deduplicator removes duplicate resources from drift reports
// before they are exported or summarized.
package deduplicator

import "github.com/snyk/driftctl-export/internal/parser"

// key uniquely identifies a resource by its type and ID.
type key struct {
	ResourceType string
	ResourceID   string
}

// Apply removes duplicate entries from all resource lists in the report.
// A duplicate is defined as a resource sharing the same Type and ResourceID.
// The first occurrence is kept; subsequent duplicates are discarded.
// The original report is mutated in place and returned for convenience.
func Apply(report *parser.Report) *parser.Report {
	if report == nil {
		return nil
	}

	report.Summary.TotalUnmanaged = len(report.UnmanagedResources)
	report.UnmanagedResources = deduplicateUnmanaged(report.UnmanagedResources)
	report.ManagedResources = deduplicateManaged(report.ManagedResources)
	report.DeletedResources = deduplicateUnmanaged(report.DeletedResources)

	return report
}

func deduplicateUnmanaged(resources []parser.Resource) []parser.Resource {
	seen := make(map[key]struct{})
	result := make([]parser.Resource, 0, len(resources))
	for _, r := range resources {
		k := key{ResourceType: r.Type, ResourceID: r.ResourceID}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, r)
	}
	return result
}

func deduplicateManaged(resources []parser.ManagedResource) []parser.ManagedResource {
	seen := make(map[key]struct{})
	result := make([]parser.ManagedResource, 0, len(resources))
	for _, r := range resources {
		k := key{ResourceType: r.Type, ResourceID: r.ResourceID}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, r)
	}
	return result
}
