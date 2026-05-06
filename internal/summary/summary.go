// Package summary provides functionality to generate human-readable
// summary statistics from parsed driftctl reports.
package summary

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/your-org/driftctl-export/internal/parser"
)

// Stats holds aggregated statistics derived from a drift report.
type Stats struct {
	TotalResources   int
	ManagedCount     int
	UnmanagedCount   int
	DeletedCount     int
	DriftedCount     int
	CoveragePercent  float64
}

// Compute calculates summary statistics from the given report.
func Compute(report *parser.Report) Stats {
	s := Stats{
		TotalResources: len(report.ManagedResources) +
			len(report.UnmanagedResources) +
			len(report.DeletedResources),
		ManagedCount:   len(report.ManagedResources),
		UnmanagedCount: len(report.UnmanagedResources),
		DeletedCount:   len(report.DeletedResources),
	}

	for _, r := range report.ManagedResources {
		if len(r.Differences) > 0 {
			s.DriftedCount++
		}
	}

	if s.TotalResources > 0 {
		s.CoveragePercent = float64(s.ManagedCount) / float64(s.TotalResources) * 100
	}

	return s
}

// Write renders the Stats table to the provided writer.
func Write(w io.Writer, s Stats) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "Metric\tCount")
	fmt.Fprintln(tw, "------\t-----")
	fmt.Fprintf(tw, "Total Resources\t%d\n", s.TotalResources)
	fmt.Fprintf(tw, "Managed\t%d\n", s.ManagedCount)
	fmt.Fprintf(tw, "Unmanaged\t%d\n", s.UnmanagedCount)
	fmt.Fprintf(tw, "Deleted\t%d\n", s.DeletedCount)
	fmt.Fprintf(tw, "Drifted (managed)\t%d\n", s.DriftedCount)
	fmt.Fprintf(tw, "Coverage\t%.1f%%\n", s.CoveragePercent)

	return tw.Flush()
}
