// Package deduplicator provides functionality to remove duplicate resource
// entries from a parsed drift report.
//
// Duplicates are identified by the combination of a resource's Type and
// ResourceID fields. When duplicates are found, the first occurrence in
// the list is retained and subsequent occurrences are discarded.
//
// This is useful when multiple glob patterns in a pipeline match the same
// underlying report file, or when a driftctl report itself contains
// repeated entries due to provider quirks.
//
// Usage:
//
//	cleaned := deduplicator.Apply(report)
package deduplicator
