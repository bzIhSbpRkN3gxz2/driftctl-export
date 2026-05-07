// Package sorter provides functionality to sort drift report resources
// by various fields such as resource type, ID, or drift status.
package sorter

import (
	"fmt"
	"sort"

	"github.com/snyk/driftctl/pkg/analyser"
)

// Field represents a sortable field on a drift resource.
type Field string

const (
	// ByType sorts resources by their resource type.
	ByType Field = "type"
	// ByID sorts resources by their resource ID.
	ByID Field = "id"
	// ByStatus sorts resources by drift status (managed, unmanaged, deleted).
	ByStatus Field = "status"
)

// Order represents the sort direction.
type Order string

const (
	// Asc sorts in ascending order.
	Asc Order = "asc"
	// Desc sorts in descending order.
	Desc Order = "desc"
)

// Options holds configuration for sorting.
type Options struct {
	Field Field
	Order Order
}

// Apply sorts the resources in the given analysis report according to opts.
// It returns an error if the field is not recognized.
func Apply(report *analyser.Analysis, opts Options) error {
	if report == nil {
		return nil
	}

	switch opts.Field {
	case ByType:
		sort.SliceStable(report.Managed(), func(i, j int) bool {
			return compareStrings(report.Managed()[i].Type, report.Managed()[j].Type, opts.Order)
		})
		sort.SliceStable(report.Unmanaged(), func(i, j int) bool {
			return compareStrings(report.Unmanaged()[i].Type, report.Unmanaged()[j].Type, opts.Order)
		})
		sort.SliceStable(report.Deleted(), func(i, j int) bool {
			return compareStrings(report.Deleted()[i].Type, report.Deleted()[j].Type, opts.Order)
		})
	case ByID:
		sort.SliceStable(report.Managed(), func(i, j int) bool {
			return compareStrings(report.Managed()[i].Id, report.Managed()[j].Id, opts.Order)
		})
		sort.SliceStable(report.Unmanaged(), func(i, j int) bool {
			return compareStrings(report.Unmanaged()[i].Id, report.Unmanaged()[j].Id, opts.Order)
		})
		sort.SliceStable(report.Deleted(), func(i, j int) bool {
			return compareStrings(report.Deleted()[i].Id, report.Deleted()[j].Id, opts.Order)
		})
	default:
		return fmt.Errorf("unsupported sort field: %q", opts.Field)
	}

	return nil
}

func compareStrings(a, b string, order Order) bool {
	if order == Desc {
		return a > b
	}
	return a < b
}
